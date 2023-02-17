/*
MIT License

Copyright (c) Microsoft Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE
*/

package memorystate

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/azure/symphony/coa/pkg/apis/v1alpha2"
	states "github.com/azure/symphony/coa/pkg/apis/v1alpha2/providers/states"
	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	Name  string
	Value int
}

func TestInitWithEmptyConfig(t *testing.T) {
	provider := MemoryStateProvider{}
	err := provider.Init(MemoryStateProviderConfig{})
	assert.Nil(t, err)
}

func TestUpSert(t *testing.T) {
	provider := MemoryStateProvider{}
	err := provider.Init(MemoryStateProvider{})
	assert.Nil(t, err)
	id, err := provider.Upsert(context.Background(), states.UpsertRequest{
		Value: states.StateEntry{
			ID: "123",
			Body: TestPayload{
				Name:  "Random name",
				Value: 12345,
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "123", id)
}

func TestList(t *testing.T) {
	provider := MemoryStateProvider{}
	err := provider.Init(MemoryStateProvider{})
	assert.Nil(t, err)
	_, err = provider.Upsert(context.Background(), states.UpsertRequest{
		Value: states.StateEntry{
			ID: "123",
			Body: TestPayload{
				Name:  "Random name",
				Value: 12345,
			},
		},
	})
	assert.Nil(t, err)
	entries, _, err := provider.List(context.Background(), states.ListRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "123", entries[0].ID)
}

func TestDelete(t *testing.T) {
	provider := MemoryStateProvider{}
	err := provider.Init(MemoryStateProvider{})
	assert.Nil(t, err)
	_, err = provider.Upsert(context.Background(), states.UpsertRequest{
		Value: states.StateEntry{
			ID: "123",
			Body: TestPayload{
				Name:  "Random name",
				Value: 12345,
			},
		},
	})
	assert.Nil(t, err)
	err = provider.Delete(context.Background(), states.DeleteRequest{
		ID: "123",
	})
	assert.Nil(t, err)
	entries, _, err := provider.List(context.Background(), states.ListRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(entries))
}

func TestMemoryStateProviderConfigFromMapNil(t *testing.T) {
	_, err := MemoryStateProviderConfigFromMap(nil)
	assert.Nil(t, err)
}

func TestMemoryStateProviderConfigFromMapEmpty(t *testing.T) {
	_, err := MemoryStateProviderConfigFromMap(map[string]string{})
	assert.Nil(t, err)
}
func TestMemoryStateProviderConfigFromMap(t *testing.T) {
	config, err := MemoryStateProviderConfigFromMap(map[string]string{
		"name": "my-name",
	})
	assert.Nil(t, err)
	assert.Equal(t, "my-name", config.Name)
}
func TestMemoryStateProviderConfigFromMapEnvOverride(t *testing.T) {
	os.Setenv("my-name", "real-name")
	config, err := MemoryStateProviderConfigFromMap(map[string]string{
		"name": "$env:my-name",
	})
	assert.Nil(t, err)
	assert.Equal(t, "real-name", config.Name)
}
func TestGet(t *testing.T) {
	provider := MemoryStateProvider{}
	err := provider.Init(MemoryStateProvider{})
	assert.Nil(t, err)
	_, err = provider.Upsert(context.Background(), states.UpsertRequest{
		Value: states.StateEntry{
			ID: "123",
			Body: TestPayload{
				Name:  "Random name",
				Value: 12345,
			},
		},
	})
	assert.Nil(t, err)
	entity, err := provider.Get(context.Background(), states.GetRequest{
		ID: "123",
	})
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, "123", entity.ID)

	payload := TestPayload{}
	data, err := json.Marshal(entity.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(data, &payload)
	assert.Nil(t, err)
	assert.Equal(t, "Random name", payload.Name)
	assert.Equal(t, 12345, payload.Value)
	entity, err = provider.Get(context.Background(), states.GetRequest{
		ID: "890",
	})
	sczErr, ok := err.(v1alpha2.COAError)
	assert.True(t, ok)
	assert.Equal(t, v1alpha2.NotFound, sczErr.State)
}

func TestUpSertEmptyID(t *testing.T) {
	provider := MemoryStateProvider{}
	err := provider.Init(MemoryStateProvider{})
	assert.Nil(t, err)
	id, err := provider.Upsert(context.Background(), states.UpsertRequest{
		Value: states.StateEntry{
			ID: "",
			Body: TestPayload{
				Name:  "Random name",
				Value: 12345,
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "", id)
}

func TestListEmptyID(t *testing.T) {
	provider := MemoryStateProvider{}
	err := provider.Init(MemoryStateProvider{})
	assert.Nil(t, err)
	_, err = provider.Upsert(context.Background(), states.UpsertRequest{
		Value: states.StateEntry{
			ID: "",
			Body: TestPayload{
				Name:  "Random name",
				Value: 12345,
			},
		},
	})
	assert.Nil(t, err)
	entries, _, err := provider.List(context.Background(), states.ListRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "", entries[0].ID)
}

func TestDeleteEmptyID(t *testing.T) {
	provider := MemoryStateProvider{}
	err := provider.Init(MemoryStateProvider{})
	assert.Nil(t, err)
	_, err = provider.Upsert(context.Background(), states.UpsertRequest{
		Value: states.StateEntry{
			ID: "",
			Body: TestPayload{
				Name:  "Random name",
				Value: 12345,
			},
		},
	})
	assert.Nil(t, err)
	err = provider.Delete(context.Background(), states.DeleteRequest{
		ID: "",
	})
	assert.Nil(t, err)
	entries, _, err := provider.List(context.Background(), states.ListRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(entries))
}

func TestGetEmptyID(t *testing.T) {
	provider := MemoryStateProvider{}
	err := provider.Init(MemoryStateProvider{})
	assert.Nil(t, err)
	_, err = provider.Upsert(context.Background(), states.UpsertRequest{
		Value: states.StateEntry{
			ID: "",
			Body: TestPayload{
				Name:  "Random name",
				Value: 12345,
			},
		},
	})
	assert.Nil(t, err)
	entity, err := provider.Get(context.Background(), states.GetRequest{
		ID: "",
	})
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, "", entity.ID)

	payload := TestPayload{}
	data, err := json.Marshal(entity.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(data, &payload)
	assert.Nil(t, err)
	assert.Equal(t, "Random name", payload.Name)
	assert.Equal(t, 12345, payload.Value)
	entity, err = provider.Get(context.Background(), states.GetRequest{
		ID: "890",
	})
	sczErr, ok := err.(v1alpha2.COAError)
	assert.True(t, ok)
	assert.Equal(t, v1alpha2.NotFound, sczErr.State)
}
func TestClone(t *testing.T) {
	provider := MemoryStateProvider{}

	p, err := provider.Clone(MemoryStateProviderConfig{
		Name: "",
	})
	assert.NotNil(t, p)
	assert.Nil(t, err)
}
