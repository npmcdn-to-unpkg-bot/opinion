'use strict';

angular.module('myApp.Playlist', ['ngRoute'])
    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/articles/list', {
            templateUrl: window.hostnametpl + 'playlist/list.html',
            controller: 'Playlist'
        });




    }])

    .controller('Playlist', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.videos = [];
        var update = function () {
            $http({
                method: 'GET',
                url: 'http://azorestv.com:6789/'
            }).then(function successCallback(response) {
                console.log(response)
                $scope.videos = response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update()




        console.log($scope.videos)
    }])