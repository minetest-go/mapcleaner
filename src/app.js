const express = require('express');

const app = express();
app.use("/node_modules", express.static('node_modules'));
app.use(express.static('public'));
app.disable('etag');

module.exports = app;
