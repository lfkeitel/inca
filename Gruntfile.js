module.exports = function(grunt) {

  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),

    jshint: {
      all: ['server/static/scripts/*.js']
    },

    uglify: {
      prod: {
        options: {
          sourceMap: true,
          sourceMapName: 'server/static/js/maps/sourcemap.map',
          banner: '/*! <%= pkg.name %> - v<%= pkg.version %> - ' +
            '<%= pkg.license %> - <%= grunt.template.today("yyyy-mm-dd") %> */'
        },
        files: {
          'server/static/js/prod.min.js': ['server/static/scripts/*.js']
        }
      }
    },

    less: {
      dev: {
        files: [{
          expand: true,
          cwd: 'server/static/styles/less',
          src: ['*.less'],
          dest: 'server/static/styles/css/',
          ext: '.css'
        }]
      }
    },

    watch: {
      scripts: {
        files: ['server/static/scripts/*.js'],
        tasks: ['newer:jshint:all']
      },
      less: {
        files: ['server/static/styles/less/*.less'],
        tasks: ['newer:less:dev']
      }
    }
  });

  grunt.event.on('watch', function(action, filepath, target) {
    grunt.log.writeln('\n' + target + ': ' + filepath + ' has ' + action);
  });

  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-less');
  grunt.loadNpmTasks('grunt-newer');

  grunt.registerTask('default', ['jshint', 'less:dev']);
  grunt.registerTask('prod', ['jshint', 'uglify', 'less:dev'])

};
