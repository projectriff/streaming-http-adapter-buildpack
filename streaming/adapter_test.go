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

package streaming_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/projectriff/streaming-http-adapter-buildpack/streaming"
)

func testAdapter(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Layers.Path, err = ioutil.TempDir("", "function-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes adapter", func() {
		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/stub-streaming-http-adapter.tgz",
			SHA256: "6b8b3e5db74a7ce366f0d3171e9f3f8a4e95aba145fa79662e306db5bda4dfde",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		a, _ := streaming.NewAdapter(dep, dc)
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = a.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())
		Expect(filepath.Join(layer.Path, "bin", "fixture-marker")).To(BeARegularFile())
	})

}
