/**
 * @param {number} The month number (1-12)
 * @param {number} The year, not zero based, required to account for leap years
 * @return {Date[]} List with date objects for each day of the month
 */
 function getDaysInMonth(month: number, year: number) {
    var date = new Date(`${month} 1, ${year}`);
    var days = [];
    while (date.getMonth() + 1 === month) {
      days.push(new Date(date));
      date.setDate(date.getDate() + 1);
    }
    return days;
  }

  export { getDaysInMonth };