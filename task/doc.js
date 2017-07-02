const gulp = require('gulp')
const shell = require('gulp-shell')

gulp.task('doc', shell.task([
  'apidoc -i ./src/server/ -o doc/',
]))
