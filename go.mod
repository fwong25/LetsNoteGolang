module webmodule

go 1.20

require github.com/lib/pq v1.10.9 // indirect

require models v1.0.0
replace models v1.0.0 => ./src/models

require db_mgmt v1.0.0
replace db_mgmt v1.0.0 => ./src/db_mgmt

require interfaces v1.0.0
replace interfaces v1.0.0 => ./src/interfaces