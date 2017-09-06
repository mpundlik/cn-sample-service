// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

//phone struct used to create customized type in Cassandra
type phone struct {
	CountryCode int    `cql:"country_code"`
	Number      string `cql:"number"`
}

//address struct used to create customized type in Cassandra. It uses pre-defined struct phone above.
type address struct {
	Street string           `cql:"street"`
	City   string           `cql:"city"`
	Zip    string           `cql:"zip"`
	Phones map[string]phone `cql:"phones"`
}

//user struct used to represent customized type in Cassandra. It uses user-defined type address.
type user struct {
	ID        string             `cql:"id" pk:"id"`
	Addresses map[string]address `cql:"addresses"`
}

//tweet struct used to represent tweet table (id text, timeline text, text text, user text, PRIMARY KEY(id)) in Cassandra.
type tweet struct {
	ID       string `cql:"id" pk:"id"`
	Timeline string `cql:"timeline"`
	Text     string `cql:"text"`
	User     string `cql:"user"`
}
