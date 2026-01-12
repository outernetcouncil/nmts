// Copyright (c) Outernet Council and Contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utilities

import (
	"fmt"
	"io"
	"os"
	"strings"

	set "github.com/deckarep/golang-set/v2"
	"github.com/samber/lo"
	"outernetcouncil.org/nmts/v1/lib/graph"
	nmtspb "outernetcouncil.org/nmts/v1/proto"
)

type EdgeSet = set.Set[*graph.Edge]

func EdgesIncidentTo(g *graph.Graph, nodeID string) EdgeSet {
	edges := NewEdgeSet()

	for _, neighbor := range g.Neighbors(nodeID) {
		edges.Append(g.Edges(nodeID, neighbor)...)
	}

	return edges
}

func OutEdgesFrom(g *graph.Graph, nodeID string) EdgeSet {
	return NewEdgeSetOf(lo.Filter(EdgesIncidentTo(g, nodeID).ToSlice(), func(e *graph.Edge, _ int) bool {
		return e != nil && e.GetA() == nodeID
	})...)
}

func InEdgesTo(g *graph.Graph, nodeID string) EdgeSet {
	return NewEdgeSetOf(lo.Filter(EdgesIncidentTo(g, nodeID).ToSlice(), func(e *graph.Edge, _ int) bool {
		return e != nil && e.GetZ() == nodeID
	})...)
}

func NewEdgeSet() EdgeSet {
	return set.NewSet[*graph.Edge]()
}

func NewEdgeSetOf(edges ...*graph.Edge) EdgeSet {
	return set.NewSet[*graph.Edge](edges...)
}

func FilterFor(edgeSet EdgeSet, kinds ...nmtspb.RK) EdgeSet {
	return NewEdgeSetOf(lo.Filter(edgeSet.ToSlice(), func(e *graph.Edge, _ int) bool {
		return lo.Contains(kinds, e.GetKind())
	})...)
}

// Get all IDs of "A" entities in inbound edges of this direct relationship EK_A -- RK --> EK_Z, given "Z" ID.
func getInIDs(g *graph.Graph, zID string, relationship nmtspb.RK) []string {
	inIDs := set.NewSet[string]()
	for edge := range FilterFor(InEdgesTo(g, zID), relationship).Iter() {
		inIDs.Add(edge.GetA())
	}
	return inIDs.ToSlice()
}

// This is useful if you just need EK_A from a direct relationship EK_A -- RK --> EK_Z, given "Z" ID.
func GetAFromZ(g *graph.Graph, zID string, relationship nmtspb.RK, isDesiredID func(string) bool) (outputID string, ok bool) {
	return lo.Find(getInIDs(g, zID, relationship), isDesiredID)
}

// Get all IDs of "Z" entities in outbound edges of this direct relationship EK_A -- RK --> EK_Z, given "A" ID.
func getOutIDs(g *graph.Graph, aID string, relationship nmtspb.RK) []string {
	outIDs := set.NewSet[string]()
	for edge := range FilterFor(OutEdgesFrom(g, aID), relationship).Iter() {
		outIDs.Add(edge.GetZ())
	}
	return outIDs.ToSlice()
}

// This is useful if you just need EK_Z from a direct relationship EK_A -- RK --> EK_Z, given "A" ID.
func GetZFromA(g *graph.Graph, aID string, relationship nmtspb.RK, isDesiredID func(string) bool) (outputID string, ok bool) {
	return lo.Find(getOutIDs(g, aID, relationship), isDesiredID)
}

// This is useful if you just need (multiple) EK_Zs from direct relationships EK_A -- RK --> EK_Z, given "A" ID.
func GetZsFromA(g *graph.Graph, aID string, relationship nmtspb.RK, isDesiredID func(id string, _ int) bool) (outputIDs []string) {
	return lo.Filter(getOutIDs(g, aID, relationship), isDesiredID)
}

type countdown struct {
	value int
}

func (c *countdown) Decrement() {
	if !c.Stopped() {
		c.value--
	}
}

func (c *countdown) Stopped() bool {
	return c == nil || c.value <= 0
}

type GraphTrace struct {
	messages []string
}

func (trace *GraphTrace) FormatForLogging() string {
	if trace == nil || len(trace.messages) == 0 {
		return ""
	}
	return strings.Join(trace.messages, "\n")
}

func (trace *GraphTrace) EmitLogMessage(writer io.Writer) {
	if trace != nil && len(trace.messages) > 0 {
		fmt.Fprintf(writer, "%s\n", trace.FormatForLogging())
	}
}

func (trace *GraphTrace) LoggingVisitFn(onVisitFn func(*graph.Graph, string)) func(*graph.Graph, string) {
	return func(g *graph.Graph, entityID string) {
		if trace != nil {
			trace.messages = append(trace.messages, fmt.Sprintf("visit: entity=%s", entityID))
		}
		onVisitFn(g, entityID)
	}
}

func (trace *GraphTrace) LoggingTraverseFn(traverseFn graph.TraverseFunc) graph.TraverseFunc {
	return func(g *graph.Graph, fromNodeID string, candidateEdge *graph.Edge) bool {
		rval := traverseFn(g, fromNodeID, candidateEdge)
		if trace != nil {
			trace.messages = append(trace.messages, fmt.Sprintf("traverse: from=%s edge={%+v} rval=%t", fromNodeID, candidateEdge.GetRelationship(), rval))
		}
		return rval
	}
}

func (trace *GraphTrace) LoggingUntilFn(untilFn func(string) bool) func(string) bool {
	return func(entityID string) bool {
		rval := untilFn(entityID)
		if trace != nil {
			trace.messages = append(trace.messages, fmt.Sprintf("until: entity=%s rval=%t", entityID, rval))
		}
		return rval
	}
}

// Traverse function that only follows relationships internal to an
// encompassing entity, like an EK_PLATFORM or EK_NETWORK_NODE.
//
// Note: because link-type Entities can RK_TRAVERSES multiple other
// link-type Entities, e.g. in a bent pipe transponder architecture,
// walk of RK_TRAVERSES is constraited to EK_INTERFACEs and EK_PORTs.
var encompassingEntityInternalRelationships = set.NewSet[nmtspb.RK](
	nmtspb.RK_RK_AGGREGATES,
	nmtspb.RK_RK_CONTAINS,
	nmtspb.RK_RK_ORIGINATES,
	nmtspb.RK_RK_SIGNAL_TRANSITS,
)

func encompassingEntityInternalTraversal(g *graph.Graph, _ string, candidateEdge *graph.Edge) bool {
	rk := candidateEdge.GetKind()
	a := g.Node(candidateEdge.GetA()).GetEntity()
	z := g.Node(candidateEdge.GetZ()).GetEntity()
	return encompassingEntityInternalRelationships.Contains(rk) ||
		(rk == nmtspb.RK_RK_TRAVERSES && z.GetEkInterface() != nil) ||
		(rk == nmtspb.RK_RK_TRAVERSES && z.GetEkPort() != nil) ||
		(rk == nmtspb.RK_RK_TERMINATES && z.GetEkDemodulator() != nil) ||
		(rk == nmtspb.RK_RK_CONTROLS && a.GetEkRouteFn() != nil)
}

// Walks only relationships "internal" to an encompassing entity of
// a given Entity Kind, stopping at an entity of the specified kind,
// if any.
func findEncompassingEntity(g *graph.Graph, isDesiredKind func(*nmtspb.Entity) bool, startingID string) string {
	startingEntity := g.Node(startingID).GetEntity()
	if startingEntity == nil {
		return ""
	}
	if isDesiredKind(startingEntity) {
		return startingID
	}

	encompassingEntityID := ""
	c := countdown{value: 4096}

	onVisitFn := func(g *graph.Graph, entityID string) {
		if e := g.Node(entityID).GetEntity(); e != nil && isDesiredKind(e) {
			encompassingEntityID = entityID
		}
	}
	shouldStopFn := func(entityID string) bool {
		if encompassingEntityID != "" {
			return true
		}
		c.Decrement()
		return c.Stopped()
	}

	var trace *GraphTrace
	// N.B.: uncomment below to emit graph traces in unit tests.
	// trace = &GraphTrace{}

	dfs := graph.DepthFirst{
		Visit:    trace.LoggingVisitFn(onVisitFn),
		Traverse: trace.LoggingTraverseFn(encompassingEntityInternalTraversal),
	}
	dfs.Walk(g, startingID, trace.LoggingUntilFn(shouldStopFn))
	trace.EmitLogMessage(os.Stderr)

	return encompassingEntityID
}

// This finds the EK_PLATFORM --RK--> x --RK--> y -->RK... -> z, where z is the startingID and RK comes from encompassingEntityInternalTraversal()
func FindEncompassingPlatform(g *graph.Graph, startingID string) string {
	return findEncompassingEntity(g, func(e *nmtspb.Entity) bool {
		return e.GetEkPlatform() != nil
	}, startingID)
}

func FindEncompassingNetworkNode(g *graph.Graph, startingID string) string {
	return findEncompassingEntity(g, func(e *nmtspb.Entity) bool {
		return e.GetEkNetworkNode() != nil
	}, startingID)
}

// Find an EK_PORT, if any, associated with physical entities that
// might form additional model element structure "attached" to the
// EK_PORT, e.g. a modem, RF chain, and/or antenna (aperture).
func FindAssociatedPort(g *graph.Graph, startingID string) string {
	startingEntity := g.Node(startingID).GetEntity()
	if startingEntity == nil {
		return ""
	}
	if startingEntity.GetEkPort() != nil {
		return startingID
	}

	portID := ""
	c := countdown{value: 4096}

	dfs := graph.DepthFirst{
		Visit: func(g *graph.Graph, entityID string) {
			if g.Node(entityID).GetEntity().GetEkPort() != nil {
				portID = entityID
			}
		},
		Traverse: func(g *graph.Graph, _ string, candidateEdge *graph.Edge) bool {
			rk := candidateEdge.GetKind()
			return rk == nmtspb.RK_RK_ORIGINATES ||
				rk == nmtspb.RK_RK_SIGNAL_TRANSITS ||
				(rk == nmtspb.RK_RK_TERMINATES && g.Node(candidateEdge.GetZ()).GetEntity().GetEkDemodulator() != nil)
		},
	}
	dfs.Walk(g, startingID, func(entityID string) bool {
		if portID != "" {
			return true
		}
		c.Decrement()
		return c.Stopped()
	})

	return portID
}

func FindBentPipeReceiverFromTransmitter(g *graph.Graph, transmitterID string) string {
	startingEntity := g.Node(transmitterID).GetEntity()
	if startingEntity == nil {
		return ""
	}
	if startingEntity.GetEkReceiver() != nil {
		return transmitterID
	}

	receiverID := ""
	c := countdown{value: 4096}

	onVisitFn := func(g *graph.Graph, entityID string) {
		if g.Node(entityID).GetEntity().GetEkReceiver() != nil {
			receiverID = entityID
		}
	}
	shouldTraverseFn := func(g *graph.Graph, _ string, candidateEdge *graph.Edge) bool {
		// Don't graph-walk past any EK_ANTENNAs.
		a := g.Node(candidateEdge.GetA()).GetEntity()
		z := g.Node(candidateEdge.GetZ()).GetEntity()
		if a.GetEkAntenna() != nil || z.GetEkAntenna() != nil {
			return false
		}
		// All entities between RX antenna and TX antenna are
		// connected via RK_SIGNAL_TRANSITS relationships.
		rk := candidateEdge.GetKind()
		return rk == nmtspb.RK_RK_SIGNAL_TRANSITS
	}
	shouldStopFn := func(entityID string) bool {
		if receiverID != "" {
			return true
		}
		c.Decrement()
		return c.Stopped()
	}

	dfs := graph.DepthFirst{
		Visit:    onVisitFn,
		Traverse: shouldTraverseFn,
	}
	dfs.Walk(g, transmitterID, shouldStopFn)

	return receiverID
}

// Walk RK_TRAVERSES relationships to collect all lower layer entities
// of the specified type(s). Does not handle RK_AGGREGATES.
// isDesiredKind should return true for any of the types in the path,
// including the starting point.
// The returned slice of IDs is only in order of highest to lowest level
// if the graph does not branch from one entity to many entities.
func GetAllTraversedEntitiesBeneath(g *graph.Graph, isDesiredKind func(*nmtspb.Entity) bool, startingID string) []string {
	startingEntity := g.Node(startingID).GetEntity()
	if startingEntity == nil || !isDesiredKind(startingEntity) {
		return []string{}
	}

	visited := set.NewSet[string](startingID)
	return getAllTraversedEntitiesBeneathHelper(g, isDesiredKind, startingID, visited)
}

func getAllTraversedEntitiesBeneathHelper(g *graph.Graph, isDesiredKind func(*nmtspb.Entity) bool, startingID string, visited set.Set[string]) []string {
	traversedEntities := []string{startingID}
	for edge := range FilterFor(OutEdgesFrom(g, startingID), nmtspb.RK_RK_TRAVERSES).Iter() {
		nextID := edge.GetZ()
		// Don't traverse if this entity has already been visited.
		if e := g.Node(nextID).GetEntity(); e != nil && isDesiredKind(e) && !visited.Contains(nextID) {
			visited.Add(nextID)
			traversedEntities = append(traversedEntities, getAllTraversedEntitiesBeneathHelper(g, isDesiredKind, nextID, visited)...)
		}
	}
	return traversedEntities
}

// Walk RK_TRAVERSES relationships to find the lowest layer entity
// of the specified type, if any. Does not handle RK_AGGREGATES.
// isDesiredKind should return true for any of the types in the path,
// including the starting point.
func findRootTraversedEntityBeneath(g *graph.Graph, isDesiredKind func(*nmtspb.Entity) bool, startingID string) string {
	entities := GetAllTraversedEntitiesBeneath(g, isDesiredKind, startingID)
	if len(entities) == 0 {
		return ""
	}
	return entities[len(entities)-1]
}

// WARNING: This will be non-deterministic in finding a root entity if
// the graph branches from one entity to many entities.
func FindRootInterfaceBeneath(g *graph.Graph, startingID string) string {
	return findRootTraversedEntityBeneath(g, func(e *nmtspb.Entity) bool {
		return e.GetEkInterface() != nil
	}, startingID)
}

// WARNING: This will be non-deterministic in finding a root entity if
// the graph branches from one entity to many entities.
func FindRootLogicalPacketLinkBeneath(g *graph.Graph, startingID string) string {
	return findRootTraversedEntityBeneath(g, func(e *nmtspb.Entity) bool {
		return e.GetEkLogicalPacketLink() != nil
	}, startingID)
}

func GetNodeIDsFromAgentID(g *graph.Graph, agentID string) []string {
	return GetZsFromA(g, agentID, nmtspb.RK_RK_CONTROLS,
		func(id string, _ int) bool {
			return g.Node(id).GetEntity().GetEkNetworkNode() != nil
		})
}

// This looks for agent given an entity with this direct relationship:
//
//	EK_SDN_AGENT --RK_CONTROLS--> (entity)
func GetAgentIDControllingThisEntity(g *graph.Graph, entityId string) (agentID string, ok bool) {
	return GetAFromZ(g, entityId, nmtspb.RK_RK_CONTROLS,
		func(id string) bool {
			return g.Node(id).GetEntity().GetEkSdnAgent() != nil
		})
}

// This gets the agent IDs controlling the modulators/demodulators, given the interfaceID and walking these paths:
//
//	EK_INTERFACE --RK_TRAVERSES--> EK_PORT --RK_ORIGINATES--> EK_MODULATOR <--RK_CONTROLS-- EK_SDN_AGENT
//	EK_INTERFACE --RK_TRAVERSES--> EK_PORT --RK_TERMINATES--> EK_DEMODULATOR <--RK_CONTROLS-- EK_SDN_AGENT
func GetAgentIDsFromNetworkInterfaceIDViaModems(g *graph.Graph, interfaceID string) (agentIDs []string) {
	modemIDSet := set.NewSet[string]()
	modemIDSet.Append(GetModulatorIDsFromInterfaceIDViaPort(g, interfaceID)...)
	modemIDSet.Append(GetDemodulatorIDsFromInterfaceIDViaPort(g, interfaceID)...)

	agentIDSet := set.NewSet[string]()
	modemIterator := modemIDSet.Iterator()
	for modemID := range modemIterator.C {
		agentID, ok := GetAgentIDControllingThisEntity(g, modemID)
		if ok {
			agentIDSet.Add(agentID)
		}
	}
	return agentIDSet.ToSlice()
}

// EK_NETWORK_NODE --RK_CONTAINS--> EK_ROUTE_FN <--RK_CONTROLS-- EK_SDN_AGENT
func GetAgentIDFromNetworkNodeIDViaRouteFn(g *graph.Graph, networkNodeID string) (agentID string, ok bool) {
	// First we get the list of potential EK_ROUTE_FNs that are contained by this node.
	// potentialRouteFns is just the list of entities with the correct relationship but not necessarily the right type.
	potentialRouteFns := set.NewSet[string]()
	for edge := range FilterFor(OutEdgesFrom(g, networkNodeID), nmtspb.RK_RK_CONTAINS).Iter() {
		potentialRouteFns.Add(edge.GetZ())
	}

	routeFns := lo.Filter(potentialRouteFns.ToSlice(), func(id string, _ int) bool {
		return g.Node(id).GetEntity().GetEkRouteFn() != nil
	})

	for _, routeFn := range routeFns {
		agentID, ok := GetAgentIDControllingThisEntity(g, routeFn)
		if ok {
			return agentID, true
		}
	}
	return "", false
}

func GetAgentIDFromInterfaceIDViaRouteFn(g *graph.Graph, interfaceID string) (agentID string, ok bool) {
	if networkNodeID, ok := GetAFromZ(g, interfaceID, nmtspb.RK_RK_CONTAINS,
		func(id string) bool {
			return g.Node(id).GetEntity().GetEkNetworkNode() != nil
		}); ok {
		return GetAgentIDFromNetworkNodeIDViaRouteFn(g, networkNodeID)
	}
	return "", false
}

func GetPortIDsFromInterfaceID(g *graph.Graph, interfaceID string) (portIDs []string) {
	// EK_INTERFACE --RK_TRAVERSES--> EK_PORT
	return GetZsFromA(g, interfaceID, nmtspb.RK_RK_TRAVERSES,
		func(id string, _ int) bool {
			return g.Node(id).GetEntity().GetEkPort() != nil
		})
}

// This gets the modulator ID(s) from the interface ID via the port(s) by walking this graph:
//
//	EK_INTERFACE --RK_TRAVERSES--> EK_PORT --RK_ORIGINATES--> EK_MODULATOR
//
// Note that an alternative is to pull directly (no graph) from the WirelessDevice field in the
// NetworkInterface that is defined inside of original recipe NetworkNode messages.
// This can return more than 1 modulatorID.
func GetModulatorIDsFromInterfaceIDViaPort(g *graph.Graph, interfaceID string) (modulatorIDs []string) {
	// EK_INTERFACE --RK_TRAVERSES--> EK_PORT
	portIDs := GetPortIDsFromInterfaceID(g, interfaceID)
	if len(portIDs) == 0 {
		return []string{}
	}

	modulatorIDSet := set.NewSet[string]()
	for _, portID := range portIDs {
		// EK_PORT --RK_ORIGINATES--> EK_MODULATOR
		modulatorID, ok := GetZFromA(g, portID, nmtspb.RK_RK_ORIGINATES,
			func(id string) bool {
				return g.Node(id).GetEntity().GetEkModulator() != nil
			})
		if ok {
			modulatorIDSet.Add(modulatorID)
		}
	}
	return modulatorIDSet.ToSlice()
}

// This gets the demodulator ID(s) from the interface ID via the port(s) by walking this graph:
//
//	EK_INTERFACE --RK_TRAVERSES--> EK_PORT --RK_TERMINATES--> EK_DEMODULATOR
//
// Note that an alternative is to pull directly (no graph) from the WirelessDevice field in the
// NetworkInterface that is defined inside of original recipe NetworkNode messages.
// This can return more than 1 demodulatorID.
func GetDemodulatorIDsFromInterfaceIDViaPort(g *graph.Graph, interfaceID string) (modemID []string) {
	// EK_INTERFACE --RK_TRAVERSES--> EK_PORT
	portIDs := GetPortIDsFromInterfaceID(g, interfaceID)
	if len(portIDs) == 0 {
		return []string{}
	}

	demodulatorIDSet := set.NewSet[string]()
	for _, portID := range portIDs {
		// EK_PORT --RK_TERMINATES--> EK_DEMODULATOR
		demodulatorID, ok := GetZFromA(g, portID, nmtspb.RK_RK_TERMINATES,
			func(id string) bool {
				return g.Node(id).GetEntity().GetEkDemodulator() != nil
			})
		if ok {
			demodulatorIDSet.Add(demodulatorID)
		}
	}
	return demodulatorIDSet.ToSlice()
}

func GetLogicalPacketLinksOriginatingFromInterface(g *graph.Graph, interfaceID string) []*nmtspb.Entity {
	logicalPacketLinkIds := GetZsFromA(g, interfaceID, nmtspb.RK_RK_ORIGINATES,
		func(id string, _ int) bool {
			return g.Node(id).GetEntity().GetEkLogicalPacketLink() != nil
		})
	logicalPacketLinks := make([]*nmtspb.Entity, len(logicalPacketLinkIds))
	for i, id := range logicalPacketLinkIds {
		logicalPacketLinks[i] = g.Node(id).GetEntity()
	}
	return logicalPacketLinks
}

// Returns the route functions associated with an interface by first, finding
// the interface's parent network node and then finding any route functions
// contained by the node.
func GetRouteFnsFromInterface(g *graph.Graph, interfaceID string) []*nmtspb.Entity {
	networkNodeID := FindEncompassingNetworkNode(g, interfaceID)
	if networkNodeID == "" {
		return []*nmtspb.Entity{}
	}

	routeFnIds := GetZsFromA(g, networkNodeID, nmtspb.RK_RK_CONTAINS,
		func(id string, _ int) bool {
			return g.Node(id).GetEntity().GetEkRouteFn() != nil
		})
	routeFns := make([]*nmtspb.Entity, len(routeFnIds))
	for i, id := range routeFnIds {
		routeFns[i] = g.Node(id).GetEntity()
	}
	return routeFns
}

// @Experimental
// This is still under development and may change at any time.
func GetAllEntityIDsUnderNetworkNode(g *graph.Graph, startingID string) []string {
	startingEntity := g.Node(startingID).GetEntity()
	if startingEntity == nil {
		return []string{}
	}
	if startingEntity.GetEkNetworkNode() == nil {
		return []string{}
	}

	entityIds := set.NewSet[string]()
	c := countdown{value: 4096}

	onVisitFn := func(g *graph.Graph, entityID string) {
		if e := g.Node(entityID).GetEntity(); e != nil {
			entityIds.Add(entityID)
		}
	}
	shouldStopFn := func(entityID string) bool {
		c.Decrement()
		return c.Stopped()
	}
	traverseFn := func(g *graph.Graph, _ string, candidateEdge *graph.Edge) bool {
		a := g.Node(candidateEdge.GetA()).GetEntity()
		return encompassingEntityInternalTraversal(g, "", candidateEdge) && a.GetEkPlatform() == nil
	}

	dfs := graph.DepthFirst{
		Visit:    onVisitFn,
		Traverse: traverseFn,
	}
	dfs.Walk(g, startingID, shouldStopFn)

	return entityIds.ToSlice()
}

// @Experimental
// This is still under development and may change at any time.
//
// Given a list of entity IDs that are affected by a fault, returns a list of all
// transitively affected entity IDs.
func ComputeTransitivelyAffectedIDsForFault(g *graph.Graph, affectedEntityIDs []string) []string {
	computedEntityIDs := computeTransitivelyAffectedIDsForFaultHelper(g, set.NewSet[string](affectedEntityIDs...), set.NewSet[string]())
	// Remove any transitively affected entities that are in affected entities.
	// This should be a no-op if none of the affectedEntityIDs are in computedEntityIds.
	computedEntityIDs.RemoveAll(affectedEntityIDs...)
	return set.Sorted(computedEntityIDs)
}

func computeTransitivelyAffectedIDsForFaultHelper(g *graph.Graph, affectedEntityIDs set.Set[string], allEntityIDs set.Set[string]) set.Set[string] {
	allEntityIDs = allEntityIDs.Union(affectedEntityIDs)
	computedEntityIDs := set.NewSet[string]()
	for entityID := range affectedEntityIDs.Iter() {
		e := g.Node(entityID).GetEntity()
		switch {
		// EK_INTERFACE: any interfaces that RK_TRAVERSES this interface.
		case e.GetEkInterface() != nil:
			computedEntityIDs.Append(getInIDs(g, entityID, nmtspb.RK_RK_TRAVERSES)...)
		// EK_LOGICAL_PACKET_LINK: any logical packet links or physical medium links that this logical packet link traverses.
		case e.GetEkLogicalPacketLink() != nil:
			computedEntityIDs.Append(GetAllTraversedEntitiesBeneath(g, func(e *nmtspb.Entity) bool {
				return e.GetEkLogicalPacketLink() != nil || e.GetEkPhysicalMediumLink() != nil
			}, entityID)...)
		// EK_PHYSICAL_MEDIUM_LINK: the logical packet link that RK_TRAVERSES this physical medium link (if any).
		case e.GetEkPhysicalMediumLink() != nil:
			if lplID, ok := GetAFromZ(g, entityID, nmtspb.RK_RK_TRAVERSES,
				func(id string) bool {
					return g.Node(id).GetEntity().GetEkLogicalPacketLink() != nil
				}); ok {
				computedEntityIDs.Add(lplID)
			}
		// EK_PLATFORM: all entities under the platform.
		case e.GetEkPlatform() != nil:
			containedIDs := getOutIDs(g, entityID, nmtspb.RK_RK_CONTAINS)
			computedEntityIDs.Append(containedIDs...)
			networkNodeIDs := lo.Filter(containedIDs,
				func(id string, _ int) bool {
					return g.Node(id).GetEntity().GetEkNetworkNode() != nil
				})
			for _, networkNodeID := range networkNodeIDs {
				computedEntityIDs.Append(GetAllEntityIDsUnderNetworkNode(g, networkNodeID)...)
			}
		// EK_NETWORK_NODE: all entities under the network node.
		case e.GetEkNetworkNode() != nil:
			computedEntityIDs.Append(GetAllEntityIDsUnderNetworkNode(g, entityID)...)
		}
		// The following entity kinds do not transitively affect any other entities:
		// EK_PORT
		// EK_ANTENNA
		// EK_MODULATOR
		// EK_DEMODULATOR
		// EK_TRANSMITTER
		// EK_RECEIVER
		// EK_ROUTE_FN
		// EK_SDN_AGENT
	}
	if computedEntityIDs.IsSubset(allEntityIDs) {
		// If this round of computation didn't add any new entities, we're done.
		return allEntityIDs
	}
	// Otherwise compute the transitively affected entities based on the new entities we found.
	computedEntityIDs.RemoveAll(allEntityIDs.ToSlice()...)
	allEntityIDs = allEntityIDs.Union(computedEntityIDs)
	return computeTransitivelyAffectedIDsForFaultHelper(g, computedEntityIDs, allEntityIDs)
}
