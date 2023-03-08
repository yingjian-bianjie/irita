package aclmapping

import (
	aclbankmapping "github.com/bianjieai/irita/aclmapping/bank"
	aclkeeper "github.com/cosmos/cosmos-sdk/x/accesscontrol/keeper"
)

type CustomDependencyGenerator struct{}

func NewCustomDependencyGenerator() CustomDependencyGenerator {
	return CustomDependencyGenerator{}
}

func (customDepGen CustomDependencyGenerator) GetCustomDependencyGenerators() aclkeeper.DependencyGeneratorMap {
	dependencyGeneratorMap := make(aclkeeper.DependencyGeneratorMap)
	dependencyGeneratorMap = dependencyGeneratorMap.Merge(aclbankmapping.GetBankDepedencyGenerator())
	return dependencyGeneratorMap
}
