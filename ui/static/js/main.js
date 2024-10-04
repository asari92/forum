var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}


    document.addEventListener('DOMContentLoaded', function() {
        document.querySelectorAll('.comment-time').forEach(function(timeElement) {
            const dateTimeString = timeElement.getAttribute('data-time'); // Получаем дату из data-time
            const dateTime = new Date(dateTimeString); // Создаем объект Date

            // Проверяем, что дата валидна
            if (!isNaN(dateTime.getTime())) {
                // Форматируем дату и время в 24-часовом формате
                const options = {
                    year: 'numeric',
                    month: 'long', // Можно использовать 'short' для сокращенного формата
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                    hour12: false, // Установите false для 24-часового формата
                };
                const formattedDate = dateTime.toLocaleString(undefined, options);
                timeElement.textContent = formattedDate; // Устанавливаем текстовое содержимое
            } else {
                timeElement.textContent = 'Invalid date'; // Сообщение об ошибке
            }
        });
    });
