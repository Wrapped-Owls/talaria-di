package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/internal/utils"
)

type (
	StdInjector struct {
		allowOverride  bool
		reflectOpts    types.ReflectionOptions
		parentInjector types.DependencyRetriever
		cacheStorage   types.Storage[types.BindKey]
	}
)

func New(canOverride bool, reflectOpts types.ReflectionOptions, parent ...types.Injector) *StdInjector {
	var parentInjector types.Injector
	if len(parent) > 0 {
		parentInjector = parent[0]
	}
	return &StdInjector{
		allowOverride:  canOverride,
		parentInjector: parentInjector,
		reflectOpts:    reflectOpts,
		cacheStorage:   NewElementsStorage[types.BindKey](canOverride, reflectOpts),
	}
}

func (s *StdInjector) SubInjector(overrides ...bool) types.Injector {
	canOverride := s.allowOverride
	if len(overrides) > 0 {
		canOverride = overrides[0]
	}

	return New(canOverride, s.reflectOpts, s)
}

func (s StdInjector) ReflectOpts() types.ReflectionOptions {
	return s.reflectOpts
}

func (s *StdInjector) Bind(key types.BindKey, value any) error {
	if s.cacheStorage.Set(key, value) {
		return utils.ErrAlreadyBound
	}
	return nil
}

func (s *StdInjector) BindNamed(bType types.BindKey, name string, value any) error {
	if s.cacheStorage.SetNamed(bType, name, value) {
		return utils.ErrAlreadyBound
	}
	return nil
}

func (s StdInjector) Get(key types.BindKey) (result any, err error) {
	if result, err = s.cacheStorage.Get(key); err != nil && s.parentInjector != nil {
		result, err = s.parentInjector.Get(key)
		if err != nil {
			err = utils.ErrNoElementFoundInsideOrParent
		}
	}
	return
}

func (s StdInjector) GetNamed(bType types.BindKey, name string) (result any, err error) {
	if result, err = s.cacheStorage.GetNamed(bType, name); err != nil && s.parentInjector != nil {
		result, err = s.parentInjector.GetNamed(bType, name)
		if err != nil {
			err = utils.ErrNoElementFoundInsideOrParent
		}
	}
	return
}

func (s *StdInjector) Set(key types.BindKey, value any) bool {
	return s.cacheStorage.Set(key, value)
}

func (s *StdInjector) SetNamed(elementType types.BindKey, name string, value any) bool {
	return s.cacheStorage.SetNamed(elementType, name, value)
}
