'use strict';

angular.module('myApp.Clients', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/clients/list', {
            templateUrl: 'clients/list.html',
            controller: 'ListClients'
        });
        $routeProvider.when('/clients/create', {
            templateUrl: 'clients/create.html',
            controller: 'CreateClients'
        });

        $routeProvider.when('/clients/edit/:id', {
            templateUrl: 'clients/edit.html',
            controller: 'EditClients'
        });
        $routeProvider.when('/tokens/create', {
            templateUrl: 'clients/createtoken.html',
            controller: 'CreateTokens'
        });

        $routeProvider.when('/tokens/list', {
            templateUrl: 'clients/listtokens.html',
            controller: 'ListTokens'
        });

        $routeProvider.when('/tokens/edit/:id', {
            templateUrl: 'clients/createtoken.html',
            controller: 'EditTokens'
        });
    }])

    .controller('ListClients', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.clients = [];
        var update = function () {
            $http({
                method: 'GET',
                url: window.hostname + 'clients/getall'
            }).then(function successCallback(response) {
                console.log(response)
                $scope.clients = response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update();
        $scope.delete = function (id) {

            $http({
                method: 'POST',
                url: window.hostname + 'clients/delete/' + id,
            }).then(function successCallback(response) {
                console.log(response);
                update();
                toastr.success('Success!', 'Client Deleted');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }


    }])
    .controller('ListTokens', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.tokens = [];
        var update = function () {
            $http({
                method: 'GET',
                url: window.hostname + 'tokens/getall'
            }).then(function successCallback(response) {
                console.log(response)
                $scope.tokens = response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update();
        $scope.delete = function (id) {

            $http({
                method: 'POST',
                url: window.hostname + 'tokens/delete/' + id,
            }).then(function successCallback(response) {
                console.log(response)
                update();
                toastr.success('Success!', 'Token Deleted');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }


    }])
    .controller('CreateClients', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.client = {};

        $scope.submit = function () {
            console.log($scope.client)
            $http({
                method: 'POST',
                url: window.hostname + 'clients/create',
                data: $scope.client
            }).then(function successCallback(response) {
                console.log(response)
                toastr.success('Success!', 'Client Created');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }])

    .controller('CreateTokens', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.token = {};
        $scope.clients=[]
        var getAllClients= function(){
            $http({
                method: 'Get',
                url: window.hostname + 'clients/getall',

            }).then(function successCallback(response) {
               $scope.clients=response.data

                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        };
        getAllClients();

        $scope.submit = function () {
            console.log($scope.client)
            $http({
                method: 'POST',
                url: window.hostname + 'tokens/create',
                data: $scope.token
            }).then(function successCallback(response) {
                console.log(response)
                toastr.success('Success!', 'Token Created');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }])
    .controller('EditClients', ['$scope', '$http', '$routeParams', 'toastr', function ($scope, $http, $routeParams, toastr) {

        var id = $routeParams.id;
        $scope.client = {};

        $http({
            method: 'GET',
            url: window.hostname + 'clients/get/' + id
        }).then(function successCallback(response) {
            console.log(response)
            $scope.client = response.data;
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response)
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });


        $scope.submit = function () {
            $http({
                method: 'POST',
                url: window.hostname + 'clients/edit/' + id,
                data: $scope.client
            }).then(function successCallback(response) {
                console.log(response)
                toastr.success('Success!', 'Publisher Edited');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }]).controller('EditTokens', ['$scope', '$http', '$routeParams', 'toastr', function ($scope, $http, $routeParams, toastr) {

    var id = $routeParams.id;
    $scope.token = {};

    $scope.clients=[]
    var getAllClients= function(){
        $http({
            method: 'Get',
            url: window.hostname + 'clients/getall',

        }).then(function successCallback(response) {
            $scope.clients=response.data

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response)
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });

    };
    getAllClients();

    $http({
        method: 'GET',
        url: window.hostname + 'tokens/get/' + id
    }).then(function successCallback(response) {
        console.log(response)
        $scope.token = response.data;
        // this callback will be called asynchronously
        // when the response is available
    }, function errorCallback(response) {
        console.log(response)
        // called asynchronously if an error occurs
        // or server returns response with an error status.
    });


    $scope.submit = function () {
        $http({
            method: 'POST',
            url: window.hostname + 'tokens/edit/' + id,
            data: $scope.token
        }).then(function successCallback(response) {
            console.log(response);
            toastr.success('Success!', 'Token Edited');
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });
    }

}]);