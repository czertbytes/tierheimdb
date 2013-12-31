module.exports = (grunt) ->
  grunt.initConfig
    pkg: grunt.file.readJSON 'package.json'

    clean:
      dist:
        files: [
          dot: true
          src: [
            '.tmp'
            'public/*'
          ]
        ]

    coffeelint:
      dist: [ 'app/scripts/{,*/}*.coffee' ]

    coffee:
      dist:
        options:
          join: true
        files:
          '.tmp/app.js': [ 'app/scripts/{,*/}*.coffee' ]

    uglify:
      dist:
        options:
          report: 'min'
          preserveComments: 'false'
        files:
          'public/scripts.min.js': [
            'app/assets/js/jquery-2.0.3.min.js'
            'app/assets/js/angular.min.js'
            'app/assets/js/angular-route.min.js'
            'app/assets/js/dropdown.js'
            'app/assets/js/infscroll.js'
            'app/assets/js/angular-google-analytics.js'
            'app/assets/js/tierheimdb.js'
            'app/assets/js/plusone.js'
            'app/assets/js/facebook.js'
            'app/assets/js/social.js'
            '.tmp/app.js'
          ]

    cssmin:
      dist:
        keepSpecialComments: 0
        files:
          'public/styles.min.css': [
            'app/assets/css/bootstrap-readable.min.css'
            'app/styles/tierheimdb.css'
            'app/styles/fadingbars.css'
          ]

    copy:
      dist:
        files: [
          expand: true
          dot: true
          flatten: true
          dest: 'public'
          src: [
            'app/*.{ico,txt,html}'
          ]
        ]
      views:
        files: [
          expand: true
          dot: true
          flatten: true
          dest: 'public/views'
          src: [
            'app/views/{,*/}*.html'
          ]
        ]
      fonts:
        files: [
          expand: true
          dot: true
          flatten: true
          dest: 'public/fonts'
          src: [
            'app/assets/fonts/{,*/}*.{eot,otf,woff,ttf}'
          ]
        ]

  grunt.loadNpmTasks 'grunt-contrib-clean'
  grunt.loadNpmTasks 'grunt-contrib-coffee'
  grunt.loadNpmTasks 'grunt-contrib-concat'
  grunt.loadNpmTasks 'grunt-contrib-copy'
  grunt.loadNpmTasks 'grunt-contrib-csslint'
  grunt.loadNpmTasks 'grunt-contrib-cssmin'
  grunt.loadNpmTasks 'grunt-contrib-htmlmin'
  grunt.loadNpmTasks 'grunt-contrib-jshint'
  grunt.loadNpmTasks 'grunt-contrib-uglify'
  grunt.loadNpmTasks 'grunt-coffeelint'

  grunt.registerTask 'build', [
    'clean:dist'
    'coffeelint:dist'
    'coffee:dist'
    'uglify:dist'
    'cssmin:dist'
    'copy:dist'
    'copy:views'
    'copy:fonts'
  ]

  grunt.registerTask 'make-js', [
    'coffeelint:dist'
    'coffee:dist'
    'uglify:dist'
  ]

  grunt.registerTask 'default', [ 'build' ]
