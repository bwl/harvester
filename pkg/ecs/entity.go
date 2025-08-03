package ecs

type Entity uint64

type entityGen uint32

type entityIndex uint32

type entityID struct {
	idx entityIndex
	gen entityGen
}
