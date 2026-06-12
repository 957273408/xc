const XLSX = require('xlsx');
const path = require('path');

const filePath = path.join(__dirname, '战队信息表-2026-06-11 16_48_20(1).xlsx');
console.log('Reading file:', filePath);

const workbook = XLSX.readFile(filePath);
console.log('Sheet names:', workbook.SheetNames);

for (const sheetName of workbook.SheetNames) {
    console.log(`\n=== Sheet: ${sheetName} ===`);
    const worksheet = workbook.Sheets[sheetName];
    const data = XLSX.utils.sheet_to_json(worksheet, { header: 1 });
    data.forEach((row, i) => {
        if (row.some(cell => cell !== undefined && cell !== null && cell !== '')) {
            console.log(`Row ${i + 1}:`, JSON.stringify(row));
        }
    });
}
