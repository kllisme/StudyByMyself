const gulp = require('gulp')

require('./task/git')
require('./task/doc')

gulp.task('release', ['git:release'])
