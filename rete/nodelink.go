package rete

import (
	"context"
	"strconv"

	"github.com/project-flogo/rules/common/model"
)

//nodelink connects 2 nodes, a rete building block
type nodeLink interface {
	String() string
	getChild() node
	isRightNode() bool

	setChild(child node)
	setIsRightChild(isRight bool)
	propagateObjects(ctx context.Context, handles []reteHandle)
}

type nodeLinkImpl struct {
	convert        []int
	numIdentifiers int

	parent    node
	parentIds []model.TupleType

	child    node
	childIds []model.TupleType

	isRight bool
	id      int
}

func newNodeLink(nw Network, parent node, child node, isRight bool) nodeLink {
	nl := nodeLinkImpl{}
	nl.initNodeLink(nw, parent, child, isRight)
	return &nl
}

func (nl *nodeLinkImpl) initNodeLink(nw Network, parent node, child node, isRight bool) {
	nl.id = nw.incrementAndGetId()
	nl.child = child
	nl.isRight = isRight

	switch v := child.(type) {

	case *joinNodeImpl:
		if isRight {
			nl.childIds = v.rightIdrs
		} else {
			nl.childIds = v.leftIdrs
		}
	case *nodeImpl:
		nl.childIds = v.identifiers
	}
	nl.parent = parent
	nl.setConvert()
	parent.addNodeLink(nl)
}

//initialize node link : for use with ClassNodeLink
func initClassNodeLink(nw Network, nl *nodeLinkImpl, child node) {
	nl.id = nw.incrementAndGetId()
	nl.child = child
	nl.childIds = child.getIdentifiers()
}

func (nl *nodeLinkImpl) getChild() node {
	return nl.child
}

func (nl *nodeLinkImpl) setConvert() {

	if len(nl.parentIds) != len(nl.childIds) {
		//TODO: ERROR handling
	}
	nl.numIdentifiers = len(nl.parentIds)
	nl.convert = make([]int, nl.numIdentifiers)

	for i := 0; i < nl.numIdentifiers; i++ {
		found := false
		for j := 0; j < nl.numIdentifiers; j++ {
			if nl.parentIds[i] == nl.childIds[j] {
				found = true
				nl.convert[i] = j
				break
			}
		}
		if !found {
			//TODO: ERROR handling
		}
	}

	need := false
	for i := 0; i < nl.numIdentifiers; i++ {
		if nl.convert[i] != i {
			need = true
			break
		}
	}
	if !need {
		nl.convert = nil
	}
}

func (nl *nodeLinkImpl) String() string {
	nextNode := ""
	switch nl.child.(type) {
	case *joinNodeImpl:
		if nl.isRight {
			nextNode += "j" + strconv.Itoa(nl.child.getID()) + "R"
		} else {
			nextNode += "j" + strconv.Itoa(nl.child.getID()) + "L"
		}
	case *filterNodeImpl:
		nextNode += "f" + strconv.Itoa(nl.child.getID())
	}
	return "link (" + nextNode + ")"
}

func (nl *nodeLinkImpl) isRightNode() bool {
	return nl.isRight
}

func (nl *nodeLinkImpl) setChild(child node) {
	nl.child = child
}
func (nl *nodeLinkImpl) setIsRightChild(isRight bool) {
	nl.isRight = isRight
}

func (nl *nodeLinkImpl) propagateObjects(ctx context.Context, handles []reteHandle) {
	if nl.convert != nil {
		convertedHandles := make([]reteHandle, nl.numIdentifiers)
		for i := 0; i < nl.numIdentifiers; i++ {
			convertedHandles[nl.convert[i]] = handles[i]
		}
		handles = convertedHandles
	}
	nl.child.assertObjects(ctx, handles, nl.isRightNode())
}
