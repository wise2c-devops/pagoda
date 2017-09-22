package main

/**
*
* @api {GET} /v1/clusters clusters retrieve
* @apiName retrieve clusters
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters
*
* @apiSuccessExample {json} Success-Response:
  [
     {
        "id": "xxx",
        "name": "xxx",
        "description": "xxx",
        "state": "initial"
     },
     {
        "id": "xxx",
        "name": "xxx",
        "description": "xxx",
        "state": "success"
     },
     {
        "id": "xxx",
        "name": "xxx",
        "description": "xxx",
        "state": "failed"
     }
  ]
*
*
*/

/**
*
* @api {GET} /v1/clusters/:id cluster retrieve
* @apiName retrieve cluster
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample {type} Request-Example:
* http://172.20.20.1:8080/v1/clusters/1
*
* @apiSuccessExample {json} Success-Response:
  {
    "id": "xxx",
    "name": "xxx",
    "description": "xxx",
    "state": "failed",
    "hosts": [
      {
        "id", "xxx",
        "hostname": "xxx",
        "ip": "xxx"
      }
    ],
    "components": [
      {
        "name": "etcd",
        "hosts": [
          "id1",
          "id2",
          "id3"
        ],
        "properties": {
          "a": "xxx",
          "b": "xxx"
        }
      },
      {
        "name": "mysql",
        "hosts": [
          "id1",
          "id2",
          "id3"
        ],
        "properties": {
          "a": "xxx",
          "b": "xxx"
        }
      }
    ]
  }
*
*
*/

/**
*
* @api {POST} /v1/clusters cluster create
* @apiName create clusters
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParam  {Object} cluster 只需要name和description字段
*
* @apiParamExample {json} Request-Example:
  {
    "name": "xxx",
    "description": "xxx"
  }
*
*
* @apiSuccessExample {json} Success-Response:
  {
    "id": "xxx",
    "name": "xxx",
    "description": "xxx",
    "state": "white"
  }
*
*
*/

/**
*
* @api {DELETE} /v1/clusters/:id cluster delete
* @apiName delete cluster
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*
*/

/**
*
* @api {PUT} /v1/clusters/:id cluster update
* @apiName update cluster
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  {
     "id": "xxx",
     "name": "xxx",
     "hosts": [
        "xxx"
     ],
     "components": [
        "etcd": {
           "hosts": [
              "xxxx"
           ]
        }
     ]
  }
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*
*/

/**
*
* @api {GET} /v1/clusters/:cluster/hosts hosts retrieve
* @apiName retrieve hosts
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/hosts
*
* @apiSuccessExample {json} Success-Response:
  HTTP/1.1 200 OK
  [
    {
      "id": "xxx",
      "hostname": "xxx",
      "ip": "xxx"
    }
  ]
*
*
*/

/**
*
* @api {POST} /v1/clusters/:cluster/hosts host create
* @apiName create host
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {json} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/hosts
	{
		"hostname": "xxx",
		"ip": "xxx"
	}
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
  {
    "id": "xxx",
		"hostname": "xxx",
		"ip": "xxx"
	}
*
*
*/

/**
*
* @api {PUT} /v1/clusters/:cluster/hosts/:host host update
* @apiName update host
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/hosts
  {
    "id": "xxx",
    "hostname": "xxx",
    "ip": "xxx"
  }
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
  {
    "id": "xxx",
    "hostname": "xxx",
    "ip": "xxx"
  }
*
*/

/**
*
* @api {DELETE} /v1/clusters/:cluster/hosts/:host host delete
* @apiName delete host
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/hosts/1
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*/

/**
*
* @api {GET} /v1/clusters/:cluster/hosts/:host host retrieve
* @apiName retrieve host
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/hosts/1
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
  {
    "id": "xxx",
    "hostname": "xxx",
    "ip": "xxx"
  }
*
*
*/

/**
*
* @api {GET} /v1/clusters/:cluster/components components retrieve
* @apiName retrieve components
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/components
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
  [
    {
      "name": "etcd"
      "hosts": [
        "id1",
        "id2"
      ],
      "properties": {
        "caFile": "xxx"
      }
    },
    {
      "name": "loadbalancer",
      "hosts": [
        "xxx",
        "xxx"
      ],
      "properties": {
        "netInterface": "eth0",
        "netMask": "24",
        "vips": [
          {
            "type": "k8s",
            "vip": "172.20.9.1"
          },
          {
            "type": "es",
            "vip": "172.20.9.2"
          },
          {
            "type": "other",
            "vip": "172.20.9.3"
          }
        ]
      }
    }
  ]
*
*
*/

/**
*
* @api {POST} /v1/clusters/:cluster/components component create
* @apiName create component
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/components
  {
    "name": "etcd",
    "hosts": [
      "id1",
      "id2"
    ]
    "properties": {}
  }
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*/

/**
*
* @api {PUT} /v1/clusters/:cluster/components/:component_id component update
* @apiName update component
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/components/idxxx
  {
    "name": "etcd",
    "hosts": [
      "xxx",
      "xxx"
    ],
    "properties": {}
  }
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
  {
    "name": "etcd",
    "hosts": [
      "xxx",
      "xxx"
    ],
    "properties": {}
  }
*
*/

/**
*
* @api {DELETE} /v1/cluster/:cluster/component/:component_id component delete
* @apiName delete component
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/component/idxxx
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*/

/**
*
* @api {GET} /v1/cluster/:cluster/component/:component_id component retrieve
* @apiName retrieve component
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {String} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/component/idxxx
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
  {
    "name": "etcd",
    "hosts": [
      "xxx",
      "xxx"
    ],
    "properties": {}
  }
*
*
*/

/**
*
* @api {PUT} /v1/clusters/:cluster_id/deployment cluster deploy/reset
* @apiName start deploy/reset cluster
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiSuccess (200) {type} name description
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/deployment
  {
    "operation": "install/reset"
  }
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*
*/

/**
*
* @api {DELETE} /v1/clusters/:cluster_id/deployment cluster stop deploy/reset
* @apiName stop deploy/reset cluster
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/deployment
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*
*/
