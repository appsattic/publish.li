// --------------------------------------------------------------------------------------------------------------------

function ajax(method, url, data, callback) {
  if ( method !== 'get' && method !== 'post' && method !== 'put' ) {
    setTimeout(function() {
      callback(new Error("Method should be get, post, or put"))
    }, 0)
  }

  if ( method === 'get' ) {
    superagent
      .get(url)
      .query(data)
      .end(function(err, res){
        if (err ) {
          console.warn(err)
          return callback('Server Error: ' + err)
        }

        var data = res.body

        console.log('err:', err)
        console.log('res:', res)
        console.log('data:', data)

        if ( !data.ok ) {
          return callback(data.msg)
        }

        callback(null, res.body.payload)
      })
    return
  }

  if ( method === 'post' || method === 'put' ) {
    superagent
      [method](url)
      .type('form') // to send as 'application/x-www-form-urlencoded'
      .send(data)
      .end(function(err, res) {
        if (err ) {
          console.warn(err)
          return callback('Server Error: ' + err)
        }

        var data = res.body

        console.log('err:', err)
        console.log('res:', res)
        console.log('data:', data)

        if ( !data.ok ) {
          return callback(data.msg)
        }

        callback(null, res.body.payload)
      })
    return
  }

}

var app = new Vue({
  el   : '#app',
  data : {
    id         : null,
    idLocal    : null,
    name       : null,
    title      : '',
    author     : '',
    website    : '',
    content    : '',
    err        : null,
    state      : 'editing',
    isEditing  : true,
    isLoading  : false,
    showSocial : false,
    twitter    : '',
    facebook   : '',
    github     : '',
    instagram  : '',
  },
  watch: {
    state : function(newState, oldState) {
      this.isEditing = newState === 'editing'
      this.isLoading = newState === 'loading'
    },
  },
  computed : {
    url : function() {
      if ( this.name ) {
        if ( typeof document.location.origin === 'undefined') {
          document.location.origin = document.location.protocol + '//' + document.location.host;
        }
        return document.location.origin + '/' + this.name
      }
      return null
    },
    saveButtonClass : function() {
      if ( this.state === 'editing' ) {
        return 'button is-success is-medium'
      }
      if ( this.state === 'loading' ) {
        return 'button is-success is-medium is-disabled is-loading'
      }
      return 'button is-success is-medium'
    },
  },
  methods : {
    onNew : function() {
      app.id = null
      app.idLocal = null
      app.name = null
      app.title = ''
      app.author = ''
      app.website = ''
      app.twitter = ''
      app.facebook = ''
      app.github = ''
      app.instagram = ''
      app.content = ''
      app.err = null
      app.state = 'editing'
    },
    onShowSocial : function() {
      app.showSocial = true
    },
    onLoad : function() {
      app.state = 'loading'
      app.err   = null

      var data = {
        id : app.idLocal,
      }
      ajax('get', '/api', data, function(err, payload) {
        // whether there is an error or not, set back to editing
        app.state = 'editing'

        if (err) {
          // stringify either an Error or a string
          app.err = err
          return
        }

        // all good, copy the data from the payload
        app.id        = payload.id
        app.idLocal   = payload.id
        app.name      = payload.name
        app.title     = payload.title
        app.author    = payload.author
        app.website   = payload.website
        app.twitter   = payload.twitter
        app.facebook  = payload.facebook
        app.github    = payload.github
        app.instagram = payload.instagram
        app.content   = payload.content
        app.err       = null
      })
    },
    onSave : function() {
      var method
      var data = {
        title     : app.title,
        author    : app.author,
        website   : app.website,
        twitter   : app.twitter,
        facebook  : app.facebook,
        github    : app.github,
        instagram : app.instagram,
        content   : app.content,
      }

      if ( app.name ) {
        // update
        method = 'post'
        data.id = app.id
        data.name = app.name
      }
      else {
        // create
        method = 'put'
      }

      // set to loading
      app.state = 'loading'
      app.err   = null

      ajax(method, '/api', data, function(err, payload) {
        // whether there is an error or not, set back to editing
        app.state = 'editing'

        if (err) {
          // stringify either an Error or a string
          app.err = err
          return
        }

        // all good, copy the data from the payload which we return on both create and update
        app.id = payload.id
        app.idLocal = payload.id
        app.name = payload.name
      })

    },
  },
})

// --------------------------------------------------------------------------------------------------------------------
