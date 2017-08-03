let glob = require("glob")
let path = require("path")

let bizPath = path.join(__dirname, "./biz/**/*.json");
let servicePath = path.join(__dirname, "./service/**/*.json");

let paths = [{
  "type": "biz",
  "path": bizPath
}, {
  "type": "service",
  "path": servicePath
}];

let status = {
  "service": {},
  "biz": {}
}

paths.forEach(function (_path) {
  let files = glob.sync(_path.path, {});
  if (files && files.length > 0) {
    files.forEach(function (file) {
      status[_path.type] = Object.assign(status[_path.type], require(file));
    })
  }
})

module.exports = status
