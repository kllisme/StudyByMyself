const package = require('../package.json')
const status = require('../resource/status')
module.exports = Object.assign({}, {
  "isDevelopment": true,
  "name": package.name,
  "defaultPassword": "e10adc3949ba59abbe56e057f20f883e",
  "server": {
    "host": "0.0.0.0",
    "port": 8081,
    "href": "https://api.erp.sodalife.dev/v1",
    "session": {
      "user": {
        "id": "USER_ID",
        "key": "USER"
      },
      "cookie": "sess",
      "expires": 3600
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
    "jwt": {
      "secret": "secret",
      "cookieName": "Authorization",
      "cookieDomain": "api.erp.sodalife.dev",
      "cookieExpire": 3600,
      "tokenExpire": 3600,
      "issuer": "api.erp.sodalife.xyz"
    },
    "cors": {
      "allowedOrigins": ["http://erp.sodalife.xyz", "http://erp.sodalife.dev", "https://erp.sodalife.xyz", "https://erp.sodalife.dev"],
      "allowedHeaders": ["Authorization", "Cookie", "Origin", "Content-Type", "*"],
      "allowedMethods": ["GET", "POST", "OPTIONS", "DELETE", "PUT"],
      "maxAge": 3600
    }
  },
  "export": {
    "loadsPath": "/temp"
  },
  "captcha": {
    "prefix": "soda:erp:api:captcha:",
    "server": "http://captcha.sodalife.xyz",
    "maxLoginRequest": 3,
    "maxSmsRequest": 3,
    "maxResetRequest": 3
  },
  "pay": {
    "remark":"苏打生活{{date}}结算款",
    "aliPay": {
      "service": {
        "batchTransNotify": "batch_trans_notify"
      },
      "partner": "",
      "inputCharset": "utf-8",
      "key": "",
      "notifyUrl": "",
      "accountName": "深圳市华策网络科技有限公司",
      "email": "laura@maizuo.com",
      "signType": "MD5",
      "requestUrl": "https://mapi.alipay.com/gateway.do"
    },
    "wechat": {
      "mchAppId": "",
      "mchId": "",
      "apiKey": "",
      "tlsFile": {
        "cert": "",
        "key": "",
        "root": ""
      },
      "checkName": "FORCE_CHECK",
      "requestUrl": {
        "createTransfers": "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers",
        "getTransfer": "https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo"
      }
    }
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
        "prefix": "soda:erp:api:",
        "maxIdle": 20,
        "maxActive": 50,
        "idleTimeout": 60
      },
      "session": {
        "addr": "192.168.1.204:6379",
        "password": "123456",
        "database": 10,
        "prefix": "soda:erp:api:session:",
        "max-idle": 20,
        "maxActive": 50,
        "idleTimeout": 60,
        "maxAgeSeconds": 3600
      },
      "rate": {
        "addr": "192.168.1.204:6379",
        "password": "123456",
        "database": 10,
        "prefix": "soda:erp:api:rate:",
        "max-retry": 3
      }
    }
  }
}, {
  package
}, {
  status
});
