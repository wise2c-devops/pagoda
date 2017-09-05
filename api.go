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
 * @api {GET} /v1/cluster/:id cluster retrieve
 * @apiName retrieve cluster
 * @apiGroup v1
 * @apiVersion  1.0.0
 *
 * @apiParamExample {type} Request-Example:
 * http://172.20.20.1:8080/v1/clusters/1
 *
 * @apiSuccessExample {json} Success-Response:
   {
	   "cluster" : {
         "id": "xxx",
         "name": "xxx",
         "description": "xxx",
         "state": "green",
         "hosts": [
            {
               "hostname": "xxx",
               "ip": "xxx"
            }
         ],
         "component": {
            "etcd": {
               "hosts": [
                  "xxx"
               ]
            },
            "loadbalancer": {
               "hosts": [
                  "xxx",
                  "xxx"
               ],
               "k8sVip": "xxx",
               "esVip": "xxx",
               "otherVip": "xxx"
            }
         }
      }
   }
 *
 *
*/

/**
*
* @api {POST} /v1/cluster cluster create
* @apiName create cluster
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
* @api {DELETE} /v1/cluster/:id cluster delete
* @apiName delete cluster
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*
*/

/**
*
* @api {PUT} /v1/cluster/:id cluster update
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
* @api {GET} /v1/cluster/:cluster/hosts hosts retrieve
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
* @api {POST} /v1/cluster/:cluster/host host create
* @apiName create host
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {json} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/host
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
* @api {PUT} /v1/cluster/:cluster/host/:host host update
* @apiName update host
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/host
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
* @api {DELETE} /v1/cluster/:cluster/host/:host host delete
* @apiName delete host
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/host/1
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
*
*/

/**
*
* @api {GET} /v1/cluster/:cluster/host/:host host retrieve
* @apiName retrieve host
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/host/1
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
* @api {GET} /v1/cluster/:cluster/components components retrieve
* @apiName retrieve components
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/components
*
*
* @apiSuccessExample {type} Success-Response:
  HTTP/1.1 200 OK
  {
    "etcd": {
      "hosts": [
        "xxx",
        "xxx"
      ]
    },
    "loadbalancer": {
      "hosts": [
        "xxx",
        "xxx"
      ]
      "k8sVip": "xxx",
      "esVip": "xxx",
      "otherVip": "xxx"
    }
  }
*
*
*/

/**
*
* @api {PUT} /v1/cluster/:cluster/component/:component component set
* @apiName set component
* @apiGroup v1
* @apiVersion  1.0.0
*
* @apiParamExample  {type} Request-Example:
  http://172.20.20.1:8080/v1/cluster/1/component/loadbalancer
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
