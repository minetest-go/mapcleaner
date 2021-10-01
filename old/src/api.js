const app = require("./app");
const state = require("./state");

app.get('/api/state', function(req, res){
	res.json(state);
});
