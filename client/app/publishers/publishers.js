'use strict';

angular.module('myApp.Publishers', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/publishers/list', {
            templateUrl: 'publishers/list.html',
            controller: 'ListPublishers'
        });
        $routeProvider.when('/publishers/create', {
            templateUrl: 'publishers/create.html',
            controller: 'CreatePublishers'
        });

        $routeProvider.when('/publishers/edit/:id', {
            templateUrl: 'publishers/edit.html',
            controller: 'EditPublishers'
        });
    }])

    .controller('ListPublishers', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.publishers = [];
        var update = function () {
            $http({
                method: 'GET',
                url: window.hostname + 'publisher/listall'
            }).then(function successCallback(response) {
                console.log(response)
                $scope.publishers = response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update()
        $scope.delete = function (id) {

            $http({
                method: 'POST',
                url: window.hostname + 'publisher/delete/' + id,
                data: $scope.publisher
            }).then(function successCallback(response) {
                console.log(response)
                update();
                toastr.success('Success!', 'Publisher Deleted');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }


    }])
    .controller('CreatePublishers', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.publisher = {};

        $scope.submit = function () {
            console.log($scope.publisher)
            $http({
                method: 'POST',
                url: window.hostname + 'publisher/create',
                data: $scope.publisher
            }).then(function successCallback(response) {
                console.log(response)
                toastr.success('Success!', 'Publisher Created');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }])
    .controller('EditPublishers', ['$scope', '$http', '$routeParams', 'toastr', function ($scope, $http, $routeParams, toastr) {

        var id = $routeParams.id;
        $scope.publisher = {};

        $http({
            method: 'GET',
            url: window.hostname + 'publisher/getid/' + id
        }).then(function successCallback(response) {
            console.log(response)
            $scope.publisher = response.data;
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
                url: window.hostname + 'publisher/edit/' + id,
                data: $scope.publisher
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

    }]);