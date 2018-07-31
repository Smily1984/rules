package model

import (
	"context"
	"time"
)

type RuleSession interface {
	AddRule(rule Rule) (int, bool)
	DeleteRule(ruleName string)

	Assert(ctx context.Context, tuple StreamTuple)
	Retract(ctx context.Context, tuple StreamTuple)
	Unregister() //remove itself from the package map
	GetName() string
	RegisterTupleDescriptors (jsonRegistry string) //a json describing types
	//NewStreamTuple(dataSource TupleTypeAlias) MutableStreamTuple
	ValidateUpdate(alias TupleTypeAlias, name string, value interface{}) bool
}

//Rule ... a Rule interface
type Rule interface {
	GetName() string
	GetIdentifiers() []TupleTypeAlias
	GetConditions() []Condition
	GetActionFn() ActionFunction
	String() string
	GetPriority() int
}

//MutableRule interface has methods to add conditions and actions
type MutableRule interface {
	Rule
	AddCondition(conditionName string, idrs []TupleTypeAlias, cFn ConditionEvaluator)
	SetAction(actionFn ActionFunction)
	SetPriority(priority int)
}

type Condition interface {
	GetName() string
	GetEvaluator() ConditionEvaluator
	GetRule() Rule
	GetIdentifiers() []TupleTypeAlias
	String() string
}

//TupleTypeAlias An internal representation of a 'DataSource'
type TupleTypeAlias string

//StreamTuple is a runtime representation of a stream of data
type StreamTuple interface {
	GetTypeAlias() TupleTypeAlias
	GetString(name string) string
	GetInt(name string) int
	GetFloat(name string) float64
	GetDateTime(name string) time.Time
	GetProperties() []string
}


//MutableStreamTuple mutable part of the stream tuple
type MutableStreamTuple interface {
	StreamTuple
	SetString(ctx context.Context, rs RuleSession, name string, value string)
	SetInt(ctx context.Context, rs RuleSession, name string, value int)
	SetFloat(ctx context.Context, rs RuleSession, name string, value float64)
	SetDatetime(ctx context.Context, rs RuleSession, name string, value time.Time)
}

//ConditionEvaluator is a function pointer for handling condition evaluations on the server side
//i.e, part of the server side API
type ConditionEvaluator func(string, string, map[TupleTypeAlias]StreamTuple) bool

//ActionFunction is a function pointer for handling action callbacks on the server side
//i.e part of the server side API
type ActionFunction func(context.Context, RuleSession, string, map[TupleTypeAlias]StreamTuple)

type ValueChangeHandler interface {
	OnValueChange(tuple StreamTuple, prop string)
}