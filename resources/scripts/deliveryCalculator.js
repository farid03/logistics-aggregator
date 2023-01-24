ymaps.ready(init);

function init() {
    // Стоимость за километр.
    let DELIVERY_TARIFF = 70,
        // Минимальная стоимость.
        MINIMUM_COST = 5000,
        myMap = new ymaps.Map('map', {
            center: [55.76, 37.64],
            zoom: 4,
            controls: []
        }),
        // Создадим панель маршрутизации.
        routePanelControl = new ymaps.control.RoutePanel({
            options: {
                // Добавим заголовок панели.
                showHeader: true,
                title: 'Расчёт доставки'
            }
        }),
        zoomControl = new ymaps.control.ZoomControl({
            options: {
                size: 'small',
                float: 'none',
                position: {
                    bottom: 145,
                    right: 10
                }
            }
        });
    // Пользователь сможет построить только автомобильный маршрут.
    routePanelControl.routePanel.options.set({
        types: {auto: true}
    });

    // Если вы хотите задать неизменяемую точку "откуда", раскомментируйте код ниже.
    /*routePanelControl.routePanel.state.set({
        fromEnabled: false,
        from: 'Москва, Льва Толстого 16'
     });*/

    myMap.controls.add(routePanelControl).add(zoomControl);

    // Получим ссылку на маршрут.
    routePanelControl.routePanel.getRouteAsync().then(function (route) {

        // Зададим максимально допустимое число маршрутов, возвращаемых мультимаршрутизатором.
        route.model.setParams({results: 1}, true);

        // Повесим обработчик на событие построения маршрута.
        route.model.events.add('requestsuccess', function () {

            let activeRoute = route.getActiveRoute();
            if (activeRoute) {
                // Получим протяженность маршрута.
                let length = route.getActiveRoute().properties.get("distance"),
                    // Вычислим стоимость доставки.
                    price = calculate(Math.round(length.value / 1000)),
                    // Создадим макет содержимого балуна маршрута.
                    balloonContentLayout = ymaps.templateLayoutFactory.createClass(
                        '<span>Расстояние: ' + length.text + '.</span><br/>' +
                        '<span style="font-weight: bold; font-style: italic">Рекомендуемая стоимость: ' + price + ' р.</span>');
                // Зададим этот макет для содержимого балуна.
                route.options.set('routeBalloonContentLayout', balloonContentLayout);
                // console.log(route.getWayPoints().get(0).properties.getAll().coordinates)
                let points = route.getWayPoints()
                let start_coords = points.get(0).geometry.getCoordinates() // start
                let end_coords = points.get(points.getLength()-1).geometry.getCoordinates() // end
                // FIXME есть неопределенное поведение с порядком начала и конца маршрута (порядок не всегда правильный)
                document.getElementById("from-lt").value = start_coords[0]
                document.getElementById("from-lg").value = start_coords[1]
                document.getElementById("to-lt").value = end_coords[0]
                document.getElementById("to-lg").value = end_coords[1]
                // Откроем балун.
                activeRoute.balloon.open();
            }
        });

    });
    // Функция, вычисляющая стоимость доставки.
    function calculate(routeLength) {
        return Math.max(routeLength * DELIVERY_TARIFF, MINIMUM_COST);
    }
}