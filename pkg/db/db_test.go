// Copyright 2020 H2O.ai, Inc.
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

package db

import (
	"strings"
	"testing"

	"github.com/h2oai/wave/pkg/keychain"
)

// Source: https://www.thegeekstuff.com/2012/09/sqlite-command-examples/

const testQueries = `
create table employee(empid integer,name varchar(20),title varchar(10));
create table department(deptid integer,name varchar(20),location varchar(10));
--
create unique index empidx on employee(empid);
--
insert into employee values(101,'John Smith','CEO');
insert into employee values(102,'Raj Reddy','Sysadmin');
insert into employee values(103,'Jason Bourne','Developer');
insert into employee values(104,'Jane Smith','Sale Manager');
insert into employee values(105,'Rita Patel','DBA');
--
insert into department values(1,'Sales','Los Angeles');
insert into department values(2,'Technology','San Jose');
insert into department values(3,'Marketing','Los Angeles');
--
select * from employee;
--
select * from department;
--
alter table department rename to dept;
alter table employee add column deptid integer;
--
update employee set deptid=3 where empid=101;
update employee set deptid=2 where empid=102;
update employee set deptid=2 where empid=103;
update employee set deptid=1 where empid=104;
update employee set deptid=2 where empid=105;
--
select * from employee;
--
create view empdept as select empid, e.name, title, d.name, location from employee e, dept d where e.deptid = d.deptid;
alter table employee add column updatedon date;
--
select * from empdept;
select empid,datetime(updatedon,'localtime') from employee;
select empid,strftime('%d-%m-%Y %w %W',updatedon) from employee;
--
drop index empidx;
drop view empdept;
drop table employee;
drop table dept;
`

var (
	testDatabaseName = "test"
)

func TestQuerying(t *testing.T) {
	kc, _ := keychain.LoadKeychain("test-keychain")
	ds := newDS(DSConf{Keychain: kc, Dir: "."})
	ds.process(DBRequest{Drop: &DropRequest{testDatabaseName}})
	batches := strings.Split(testQueries, "--")
	for _, batch := range batches {
		queries := strings.Split(batch, ";")
		var stmts []Stmt

		for _, query := range queries {
			query = strings.TrimSpace(query)
			if len(query) > 0 {
				stmts = append(stmts, Stmt{query, nil})
			}
		}

		result := ds.process(DBRequest{Exec: &ExecRequest{testDatabaseName, stmts, 1}})
		t.Log("batch", stmts)
		if reply, ok := result.(ExecReply); ok {
			if len(reply.Error) > 0 {
				t.Error(reply.Error)
			} else {
				t.Log("result", reply.Results)
			}
		} else {
			t.Error("unexpected result")
		}
	}
	ds.process(DBRequest{Drop: &DropRequest{testDatabaseName}})
}
