// export function sendPostMove(treeId): number {
//     const xhr = new XMLHttpRequest();
//     const url = 'http://localhost:8080/main'; // Замените на ваш URL

//     const requestBody = {
//         treeId: treeId,
//         positionAfter: positionAfter,
//         preventMoveId: preventMoveId,
//         colorWhite: colorWhite
//     }

//     // Создаем объект с данными для отправки
//     const data = JSON.stringify(requestBody);

//     // Инициализируем запрос
//     xhr.open('POST', url, false); // false делает запрос синхронным


//     // Устанавливаем заголовок для токена
//     //xhr.setRequestHeader('token', token); // Используем 'token' вместо 'Authorization'

//     // Устанавливаем заголовок для отправки JSON
//     xhr.setRequestHeader('Content-Type', 'application/json');

//     // Отправляем запрос
//     xhr.send(data);

    
//     // Парсим ответ в JSON
//     const response = JSON.parse(xhr.responseText);
//     // Возвращаем значение поля id

//     //console.log('aaaaaaaaaaa');
//     //console.log(response.id);
//     console.log("asasasas", response.id)
//     return response.id; // Предполагается, что id - это число
    
// }
