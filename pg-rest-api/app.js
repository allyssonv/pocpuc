'use strict'
const pg = require('pg')
const express = require('express')
const app = express()


const connection = new pg.Client({
  host: 'localhost',
  port: 5432,
  user: 'postgres',
  password: 'jarvis',
  database: 'products'
})

connection.connect()

app.use(function(req, res, next) {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});

app.get('/', (req, res) => {
  res.end('<h1>Hello</h1>')
})


app.get('/products', (req, res) => {
  connection.query('select * from products', (err, result) => {
  res.json(result)
  //connection.end()
})

})

app.listen(3001, () => {
  console.log('Up! Listening on port 3001')
})
