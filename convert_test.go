// Copyright 2023 Google LLC
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

package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/pprof/profile"
)

func TestComments(t *testing.T) {
	p, err := ConvertCSVToPprof(strings.NewReader("stack,samples/count,time/ms\nfoo;bar,1,1000"), ";")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	wantComments := []string{"Generated by csv2pprof"}
	if diff := cmp.Diff(wantComments, p.Comments); diff != "" {
		t.Errorf("wanted comments %v got %v, diff: %v", wantComments, p.Comments, diff)
	}
}

func TestUnits(t *testing.T) {
	type test struct {
		input []string
		want  []*profile.ValueType
	}

	tests := []test{
		{
			// Two units fully specified
			input: []string{
				"stack,samples/count,time/ms",
				"foo;bar,1,1000",
			},
			want: []*profile.ValueType{
				{
					Type: "samples",
					Unit: "count",
				},
				{
					Type: "time",
					Unit: "ms",
				},
			},
		},
		{
			// One unit fully specified
			input: []string{
				"stack,time/ms",
				"foo;bar,1",
			},
			want: []*profile.ValueType{
				{
					Type: "time",
					Unit: "ms",
				},
			},
		},
		{
			// No units given, default to 'count'
			input: []string{
				"stack,samples",
				"foo;bar,1",
			},
			want: []*profile.ValueType{
				{
					Type: "samples",
					Unit: "count",
				},
			},
		},
		{
			// If multiple possible units, choose the last one.
			input: []string{
				"stack,samples/unit1/unit2",
				"foo;bar,1",
			},
			want: []*profile.ValueType{
				{
					Type: "samples/unit1",
					Unit: "unit2",
				},
			},
		},
		{
			// Stack at the end
			input: []string{
				"time/seconds,stack",
				"1,foo;bar",
			},
			want: []*profile.ValueType{
				{
					Type: "time",
					Unit: "seconds",
				},
			},
		},
		{
			// Stack in the middle
			input: []string{
				"time/seconds,stack,age/years",
				"1,foo;bar,18",
			},
			want: []*profile.ValueType{
				{
					Type: "time",
					Unit: "seconds",
				},
				{
					Type: "age",
					Unit: "years",
				},
			},
		},
	}

	for i, c := range tests {
		p, err := ConvertCSVToPprof(strings.NewReader(strings.Join(c.input, "\n")), ";")
		if err != nil {
			t.Fatalf("got error: %v", err)
		}

		opts := cmpopts.IgnoreUnexported(profile.ValueType{})
		got := p.SampleType
		if diff := cmp.Diff(c.want, got, opts); diff != "" {
			t.Errorf("test %v, wanted SampleType %#v got %#v, diff: %v", i, c.want, got, diff)
		}
	}
}

func TestSamples(t *testing.T) {
	type test struct {
		input    []string
		stackSep string
		want     []*profile.Sample
	}

	tests := []test{
		{
			// Stack at the start
			input: []string{
				"stack,samples/count",
				"foo;bar,1",
			},
			stackSep: ";",
			want: []*profile.Sample{
				{Value: []int64{1}},
			},
		},
		{
			// Stack at the end
			input: []string{
				"samples/count,stack",
				"1,foo;bar",
			},
			stackSep: ";",
			want: []*profile.Sample{
				{Value: []int64{1}},
			},
		},
		{
			// Stack in the middle end
			input: []string{
				"samples/count,stack,age/years",
				"1,foo;bar,18",
			},
			stackSep: ";",
			want: []*profile.Sample{
				{Value: []int64{1, 18}},
			},
		},
		{
			input: []string{
				"stack,samples/count,time/ms",
				"foo;bar,1,1000",
			},
			stackSep: ";",
			want: []*profile.Sample{
				{Value: []int64{1, 1000}},
			},
		},
		{
			input: []string{
				"stack,samples/count,time/ms",
				"foo;bar,1,1000",
				"foo,2,2000",
			},
			stackSep: ";",
			want: []*profile.Sample{
				{Value: []int64{1, 1000}},
				{Value: []int64{2, 2000}},
			},
		},
		{
			// Stack at the start
			input: []string{
				"stack,samples/count",
				"\"foo\nbar\",1",
			},
			stackSep: "\n",
			want: []*profile.Sample{
				{Value: []int64{1}},
			},
		},
	}

	for _, c := range tests {
		p, err := ConvertCSVToPprof(strings.NewReader(strings.Join(c.input, "\n")), c.stackSep)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}

		if len(p.Sample) != len(c.want) {
			t.Errorf("wanted %v samples got %v samples. samples: %v", len(c.want), len(c.want), c.want)
			continue
		}
		got := p.Sample
		opts := []cmp.Option{
			cmpopts.IgnoreUnexported(profile.Sample{}),
			cmpopts.IgnoreFields(profile.Sample{}, "Location"),
		}
		if diff := cmp.Diff(c.want, got, opts...); diff != "" {
			t.Errorf("wanted Sample %#v got %#v diff: %v", c.want, got, diff)
		}
	}
}

func TestErrors(t *testing.T) {
	type test struct {
		input   []string
		wantErr string
	}

	tests := []test{
		{
			input: []string{
				"samples/count",
				"1",
			},
			wantErr: "expected \"stack\" in CSV header row, got: [\"samples/count\"]",
		},
		{
			input: []string{
				"stack",
				"foo;bar",
			},
			wantErr: "expected columns with weights in CSV header row, got [\"stack\"]",
		},
		{
			input: []string{
				"stack,weight",
				"foo;bar",
			},
			wantErr: "error reading CSV: record on line 2: wrong number of fields",
		},
		{
			input: []string{
				"stack,weight",
				"foo;bar,not-a-number",
			},
			wantErr: "on line 2, couldn't parse number: strconv.ParseInt: parsing \"not-a-number\": invalid syntax",
		},
	}
	for i, c := range tests {
		_, err := ConvertCSVToPprof(strings.NewReader(strings.Join(c.input, "\n")), ";")
		if err == nil {
			t.Errorf("test %v, wanted error %q, got error nil", i, c.wantErr)
			continue
		}
		if err.Error() != c.wantErr {
			t.Errorf("test %v, wanted error %q, got error: %q", i, c.wantErr, err)
		}
	}
}
