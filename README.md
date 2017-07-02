## 苏打生活 ERP API

### 目录结构

- src/static 

前端业务逻辑处理 & 页面静态资源

- src/server

后端业务逻辑处理

### 编译和运行

#### 安装golang环境

推荐下载地址

[Windows](http://golangtc.com/static/go/1.5.2/go1.5.2.windows-amd64.zip)

[Mac](http://golangtc.com/static/go/1.5.2/go1.5.2.darwin-amd64.tar.gz)


#### 安装依赖包
```
npm i 
```


#### 开始开发

- 前端即时编译
```
gulp 
```

- 启动服务端
```
make dev
```


### 发布
- 打标签并上传
	```
	git tag 20161120.1
	git push origin master:master --tags 
	```

- 登录发布环境
	```
	ssh ${deployment}
	```

- 手动设置环境
	```
	export NODE_ENV=${production, intranet}
	kubectl config set-context ${prod, dev}
	docker login ${reg.maizuo.com, reg.miz.so}
	```

- 更新发布环境的仓库
	```
	cd /data/deploy/project/golang/src/maizuo.com/soda-h5
	git fetch -p --tags
	git checkout -B master origin/master
	```

- 第一次发布
	```
	cd /data/deploy/chamber/soda-h5
	make HOSTING=${intranet,production}
	```

- 滚动更新
	```
	cd /data/deploy/chamber/soda-h5
	make update HOSTING=${intranet,production}
	```
