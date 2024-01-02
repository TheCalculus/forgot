#![allow(dead_code)]
#![allow(unused_imports)]
#![allow(unused_mut)]
#![allow(unused_variables)]

extern crate chrono;

use chrono::{DateTime, Utc};

use std::{
    collections::HashMap,
    fs::{File, OpenOptions},
    io::Write
};

struct Table<T: std::fmt::Debug + Clone> {
    creation:   DateTime<Utc>,
    updated:    DateTime<Utc>,
    aof_log:    File,
    data:       HashMap<String, Data<T>>,
    identifier: u64,
}

enum QueryError {
    ReadSuccess,
    ReadFailure,
    WriteSuccess,
    WriteFailure,
    KeyAlreadyExists,
    NotImplemented,
}

#[derive(Debug, Clone)]
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

impl<T: std::fmt::Debug + Clone> Table<T> {
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
        QueryError::WriteSuccess
    }

    fn remove(&mut self, key: &str) -> QueryError {
        self.data.remove(key);
        writeln!(self.aof_log, "remove({})", key).unwrap_or(());
        QueryError::WriteSuccess
    }
}

fn main() {
    let mut table: Table<i32> = Table::new();
    let data: Data<i32> = Data::from(0);

    table.write("nice_key".to_string(), data);
    table.remove("nice_key");
}