const Promise = require('bluebird')
const gulp = require('gulp')
const sequence = require('run-sequence')
const bump = require('gulp-bump')
const gutil = require('gulp-util')
const git = require('gulp-git')
const fs = require('fs')
const config = require('config')
const moment = require('moment')

/**
 *  format:
 *
 *    YYYYMMDD.X
 *
 *  e.g.:
 *
 *    20161121.1
 *
 */
function getVersion() {
  const today = moment().format('YYYYMMDD')
  const prefix = today + '.'
  Promise.promisify(git.exec)({ args: `tag -l "${prefix}*"` })
    .then(function (stdout) {
      const tags = stdout.trim().split('\n').map(function (tag) {
        return +tag.slice(prefix.length)
      })
      return Math.max.apply(Math, tags)
    }, function (err) {
      return 0
    })
    .then(function (latest) {
      return prefix + (latest + 1)
    })
}

gulp.task('git:release:commit', function () {
  return gulp.src('.')
    .pipe(git.commit('[Release] Bumped ' + getVersion(), {args: '-a'}))
})

gulp.task('git:release:push', function (cb) {
  git.push('origin', 'master', cb)
})

gulp.task('git:release:tag', function (cb) {
  const version = getVersion()
  git.tag(version, 'Created Tag for version: ' + version, function (error) {
    if (error) {
      return cb(error)
    }
    git.push('origin', 'master', {args: '--tags'}, cb)
  })
})

gulp.task('git:release', function (callback) {
  sequence(
    'git:release:commit',
    'git:release:push',
    'git:release:tag',
    function (error) {
      if (error) {
        console.log(error.message)
      } else {
        console.log('RELEASE FINISHED SUCCESSFULLY')
      }
      callback(error)
    })
})
