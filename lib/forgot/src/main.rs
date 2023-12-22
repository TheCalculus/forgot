#![allow(dead_code)]
#![allow(unused_imports)]
#![allow(unused_mut)]
#![allow(unused_variables)]

extern crate chrono;

use std::collections::HashMap;
use std::any::Any;
use chrono::{DateTime, Utc};
use std::fs::{File, OpenOptions};
use std::io::{self, Write};

struct Table<'a> {
    creation:   DateTime<Utc>,
    updated:    DateTime<Utc>,
    rdb_log:    File,
    identifier: u64,
    data: HashMap<&'a str, &'a Data<'a>>,
}

enum QueryError {
    ReadSuccess,
    ReadFailure,
    WriteSuccess,
    WriteFailure,
    KeyAlreadyExists,
    NotImplemented,
}

impl<'a> Table<'a> {
    fn new() -> Self {
        let mut options = OpenOptions::new();

        options.append(true);

        let file = options.open("google");

        Self {
            creation:   Utc::now(),
            updated:    Utc::now(),
            rdb_log:    file.unwrap(),
            identifier: 0,
            data:       HashMap::new(),
        }
    }

    fn write(&mut self, key: &'a str, val: &'a Data<'a>) -> QueryError {
        if !self.data.contains_key(key) {
            return QueryError::KeyAlreadyExists
        }

        self.data.insert(key, val);

        return QueryError::WriteSuccess;
    }

    fn remove(&mut self, key: &'a str, option_val: Option<&'a Data<'a>>) -> QueryError {
        if let Some(val) = option_val {
            // implement this later
            return QueryError::NotImplemented;
        }

        self.data.remove(key);

        return QueryError::WriteSuccess;
    }
}

struct Data<'a> {
    creation:   DateTime<Utc>,
    updated:    DateTime<Utc>,
    identifier: u64,
    data:       &'a dyn Any,
}

fn main() {
    let mut table = Table::new();
}
