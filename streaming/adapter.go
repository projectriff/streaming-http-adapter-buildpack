/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package streaming

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
)

type Adapter struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewAdapter(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, plan *libcnb.BuildpackPlan) Adapter {
	return Adapter{LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan)}
}

func (a Adapter) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	a.LayerContributor.Logger = a.Logger

	return a.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		a.Logger.Bodyf("Expanding to %s", layer.Path)

		file := filepath.Join(layer.Path, "bin")
		if err := os.MkdirAll(file, 0755); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to create directory %s\n%w", file, err)
		}

		if err := crush.ExtractTarGz(artifact, file, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to extract %s\n%w", artifact.Name(), err)
		}

		layer.Launch = true
		return layer, nil
	})
}

func (Adapter) Name() string {
	return "adapter"
}
