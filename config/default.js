const package = require('../package.json')

module.exports = Object.assign({}, {package}, {
  "isDevelopment": true,
  "name":package.name,
  "defaultPassword":"e10adc3949ba59abbe56e057f20f883e",
  "server": {
    "host": "0.0.0.0",
    "port": 8081,
    "href":"https://m.sodalife.xyz/v1",
    "session": {
      "user": {
        "id": "USER_ID",
        "key": "USER"
      },
      "cookie": "sess",
      "expires": 24
    },
    "static": {
      "dir": "./build/static",
    },
    "template": {
      "dir": "./src/template",
      "extname": ".html"
    },
    "favicon": {
      "path": "./build/static/vendor/soda/favicon.ico"
    },
    "log": {
      "path": "./log/log.log"
    },
    "rate": {
      "value": "10-M"
    },
    "captcha": {
      "key": "CAPTCHA",
      "font": {
        "path": "./resource/font/comic.ttf"
      }
    },
    "jwt":{
      "secret":"secret",
      "cookieName":"Authorization",
      "cookieDomain":"api.erp.sodalife.xyz",
      "cookieExpire":24,
      "tokenExpire":24,
      "issuer":"api.erp.sodalife.xyz"
    },
    "cors": {
      "allowedOrigins": ["http://erp.sodalife.xyz", "http://erp.sodalife.dev","https://erp.sodalife.xyz", "https://erp.sodalife.dev"],
      "allowedHeaders": ["Authorization","Cookie","Origin","Content-Type", "*"],
      "allowedMethods": ["GET", "POST", "OPTIONS", "DELETE", "PUT"],
      "maxAge": 3600
    }
  },
  "captcha": {
    "prefix": "erp-api:captcha:",
    "server": "http://captcha.sodalife.xyz",
    "maxLoginRequest": 3,
    "maxSmsRequest": 3,
    "maxResetRequest": 3
  },
  "resource": {
    "database": {
      "soda-manager": {
        "r": {
          "dialect": "mysql",
          "host": "192.168.1.204",
          "port": 3306,
          "user": "web",
          "password": "123456",
          "database": "soda-manager",
          "maxIdle": 20,
          "maxOpen": 20
        },
        "wr": {
          "dialect": "mysql",
          "host": "192.168.1.204",
          "port": 3306,
          "user": "web",
          "password": "123456",
          "database": "soda-manager",
          "maxIdle": 20,
          "maxOpen": 20
        }
      },
      "soda": {
        "r": {
          "dialect": "mysql",
          "host": "192.168.1.204",
          "port": 3306,
          "user": "web",
          "password": "123456",
          "database": "soda",
          "maxIdle": 20,
          "maxOpen": 20
        },
        "wr": {
          "dialect": "mysql",
          "host": "192.168.1.204",
          "port": 3306,
          "user": "web",
          "password": "123456",
          "database": "soda",
          "maxIdle": 20,
          "maxOpen": 20
        }
      }
    },
    "redis": {
      "default": {
        "addr": "192.168.1.204:6379",
        "password": "123456",
        "database": 10,
        "prefix": "soda-api:",
        "maxIdle": 20,
        "maxActive": 50,
        "idleTimeout": 60
      },
      "session": {
        "addr": "192.168.1.204:6379",
        "password": "123456",
        "database": 10,
        "prefix": "soda-api:session:",
        "max-idle": 20,
        "maxActive": 50,
        "idleTimeout": 60,
        "maxAgeSeconds": 3600,
        "expires": 24
      },
      "rate": {
        "addr": "192.168.1.204:6379",
        "password": "123456",
        "database": 10,
        "prefix": "soda-api:rate:",
        "max-retry": 3
      }
    }
  }
});
