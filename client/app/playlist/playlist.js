'use strict';

angular.module('myApp.Playlist', ['ngRoute'])
    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/playlist/list', {
            templateUrl: window.hostnametpl + 'playlist/list.html',
            controller: 'Playlist'
        });


        $routeProvider.when('/playlist/settings', {
            templateUrl: window.hostnametpl + 'playlist/settings.html',
            controller: 'Settings'
        });


    }])

    .controller('Playlist', ['$scope', '$http',  function ($scope, $http) {
        $scope.videos = [];
        $scope.Date;
        $scope.StartTime = '';
        var startv2;

        $scope.format = function (seconds) {
            return new Date(seconds * 1000).toISOString().substr(11, 8);
        };

        var update = function () {
            $http({
                method: 'GET',
                url: 'http://azorestv.com:6789/getplaylist'
            }).then(function successCallback(response) {
                console.log(response);
                $scope.videos = [];
                $scope.StartTime =  response.data.StartTime;
                console.log(response.data.StartTime);
                startv2=  new Date($scope.StartTime);



                angular.forEach(response.data.Videos, function (value, key) {
                    var res = calcVideoTime(value.Duration);
                    value.PlayingInterval = res[0];
                    value.Playing = dateCheck(res[1], res[2], new Date());
                    this.push(value);
                }, $scope.videos);

                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update();

        var startHour = 18;

        var sumseconds = 0;

        function dateCheck(from, to, check) {
            if ((check.getTime() <= to.getTime() && check.getTime() >= from.getTime())) {
                return true;
            }
            return false;
        }


        function datetoHHMMSS(date) {
            return (date.getHours() < 10 ? "0" + date.getHours() : date.getHours()) +
                ":" + (date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes()) +
                ":" + (date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds());
        }

        var calcVideoTime = function (seconds) {
            var intseconds = seconds * 1;

            var datestart = new Date(startv2.getTime() + (sumseconds * 1000));

            var dateend = new Date(datestart.getTime() + (seconds) * 1000);
            sumseconds = (sumseconds + intseconds);
            var res = [];
            res[0] = datetoHHMMSS(datestart) + ' - ' + datetoHHMMSS(dateend);
            res[1] = datestart;
            res[2] = dateend;
            return res
        };


    }]).controller('Settings', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {


    $scope.StartTime='';
    var update = function () {
        $http({
            method: 'GET',
            url: 'http://azorestv.com:6789/starttime'
        }).then(function successCallback(response) {
            console.log(response)
            $scope.StartTime=string2date(response.data)
            console.log(string2date(response.data))
            toastr.success('Success!', 'Hora de inicio alterada');
           // $scope.StartTime=response.data.StartTime

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });
    };
    update();

    function pad(n, width, z) {
        z = z || '0';
        n = n + '';
        return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
    }

    function date2string(date){
        return pad(date.getHours(),2)+":"+pad(date.getMinutes(),2)
    }

    $scope.reload= function (){
        console.log("ups")
        $http({
            method: 'POST',
            url: 'http://azorestv.com:6789/reload',
        }).then(function successCallback(response) {
            toastr.success('Success!', 'servidor reiniciado');
            console.log(response);
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });

    }

    function string2date(strdate){

        if (strdate.indexOf(":") > -1){
            var date = strdate.split(":");
            var   hour,minute;
            hour=date[0];
            minute=date[1];
           var d = new Date();
            d.setHours(hour);
            d.setMinutes(minute);
            return d
        }
        return new Date();


    }


    $scope.updateStartTime= function(start){
        console.log(date2string(start));
        $http({
            method: 'POST',
            url: 'http://azorestv.com:6789/starttime',
            data:{StartTime:date2string($scope.StartTime)},
        }).then(function successCallback(response) {
            console.log(response);

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });

    };




}]);