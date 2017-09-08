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
        "state": "green"
     },
     {
        "id": "xxx",
        "name": "xxx",
        "description": "xxx",
        "state": "red"
     },
     {
        "id": "xxx",
        "name": "xxx",
        "description": "xxx",
        "state": "white"
     }
  ]
*
*
*/

/**
*
* @api {GET} /v1/clusters/:id cluster retrieve
* @apiName retrieve clusters
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
    "state": "green",
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
          "xxx",
          "xxx",
          "xxx"
        ],
        "property": {
          "a": "xxx",
          "b": "xxx"
        }
      },
      {
        "name": "loadbalancer",
        "hosts": [
          "xxx",
          "xxx",
          "xxx"
        ],
        "property": {
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
        "xxx",
        "xxx"
      ],
      "property": {
        "caFile": "xxx"
      }
    },
    {
      "name": "loadbalancer",
      "hosts": [
        "xxx",
        "xxx"
      ],
      "property": {
        "k8sVip": "xxx",
        "esVip": "xxx",
        "otherVip": "xxx"
      }
    }
  ]
*
*
*/

/**
*
* @api {PUT} /v1/clusters/:cluster/components/:component component set
* @apiName set component
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/clusters/1/components/loadbalancer
  "loadbalancer": {
    "hosts": [
      "xxx",
      "xxx"
    ]
    "k8sVip": "xxx",
    "esVip": "xxx",
    "otherVip": "xxx"
  }
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*/

/**
*
* @api {DELETE} /v1/cluster/:cluster/component/:component component delete
* @apiName delete component
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/component/loadbalancer
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*/

/**
*
* @api {GET} /v1/cluster/:cluster/component/:component component retrieve
* @apiName retrieve component
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {String} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/component/loadbalancer
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
  {
    "hosts": [
      "xxx",
      "xxx"
    ],
    "k8sVip": "xxx",
    "esVip": "xxx",
    "otherVip": "xxx"
  }
*
*
*/
