#![allow(dead_code)]
#![allow(unused_imports)]
#![allow(unused_mut)]
#![allow(unused_variables)]

extern crate chrono;

use chrono::{DateTime, Utc};

use std::{
    collections::HashMap,
    fs::{File, OpenOptions},
    io::Write,
    option,
};

struct Table<T: std::fmt::Debug + Clone> {
    creation:   DateTime<Utc>,
    updated:    DateTime<Utc>,
    aof_log:    File,
    data:       HashMap<String, Data<T>>,
    identifier: u64,
}

enum QueryError {
    Success,
    Failure,
    Found,
    NotFound,
    KeyAlreadyExists,
    NotImplemented, 
}

/* should receive only 
    &str, Data<T> */
enum QueryType<T, U> {
    ByKey(T),
    ByValue(U),
}

#[derive(Debug, Clone, PartialEq)]
struct Data<T: Clone> {
    creation:   DateTime<Utc>,
    updated:    DateTime<Utc>,
    data:       Box<T>,
    identifier: u64,
}

impl<T: Clone> Data<T> {
    fn from(data: T) -> Self {
        let time = Utc::now();
        
        Self {
            creation:  time,
            updated:   time,
            data:      Box::new(data),
            identifier: 0,
        }
    }
}

impl<T: std::fmt::Debug + std::clone::Clone + std::cmp::PartialEq> Table<T> {
    fn new() -> Self {
        let time = Utc::now();
        let file = OpenOptions::new()
            .append(true)
            .create(true)
            .open(time.format("%Y-%m-%dT%H:%M:%SZ").to_string())
            .expect("failed to open file");

        Self {
            creation:   time,
            updated:    time,
            aof_log:    file,
            data:       HashMap::new(),
            identifier: 0,
        }
    }

    fn write(&mut self, key: String, val: Data<T>) -> QueryError {
        if self.data.contains_key(&key) {
            return QueryError::KeyAlreadyExists;
        }
        self.data.insert(key.clone(), val.clone());
        writeln!(self.aof_log, "write({}, {:?})", key, val).unwrap_or(());
        QueryError::Success
    }

    fn remove(&mut self, key: &str) -> QueryError {
        self.data.remove(key);
        writeln!(self.aof_log, "remove({})", key).unwrap_or(());
        QueryError::Success
    }

    fn get(&mut self, key: String) -> Result<&Data<T>, QueryError> {
        match self.data.get(&key) {
            Some(v) => Ok(&v),
            None    => Err(QueryError::NotFound),
        }
    }

    fn query(&mut self, key: QueryType<String, Data<T>>) -> Result<&Data<T>, QueryError> {
        match key {
            QueryType::ByKey(k) => {
                // shouldn't use query for fetching keys, but implementing anyway
                self.get(k)
            },
            QueryType::ByValue(v) => {
                for (key, value) in &self.data {
                    if *value != v { continue; }
                    return Ok(&self.data[key]);
                }
                Err(QueryError::NotFound)
            }
        }
    }
}

fn main() {
    let mut table: Table<&str> = Table::new();
    let data: Data<&str> = Data::from("hope its a good one");

    table.write("happy_new_year".to_string(), data);
}