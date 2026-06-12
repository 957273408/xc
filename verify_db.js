const mysql = require('mysql2/promise');

const DB_CONFIG = {
    host: '192.168.31.23',
    port: 13306,
    user: 'root',
    password: 'Sun123456',
    database: 'gva',
    charset: 'utf8mb4'
};

async function verifyData() {
    let connection;

    try {
        connection = await mysql.createConnection(DB_CONFIG);
        console.log('========== 数据库数据验证 ==========\n');

        // 检查战队
        const [teams] = await connection.execute('SELECT id, team_name, total_bounty FROM exa_teams ORDER BY id');
        console.log(`战队总数: ${teams.length}`);
        teams.forEach(t => {
            console.log(`  ID ${t.id}: ${t.team_name} (赏金: ${t.total_bounty})`);
        });

        // 检查选手
        const [players] = await connection.execute(`
            SELECT p.id, p.player_name, p.uid, p.team_id, t.team_name, p.bounty
            FROM exa_players p
            LEFT JOIN exa_teams t ON p.team_id = t.id
            ORDER BY p.team_id, p.id
        `);

        console.log(`\n选手总数: ${players.length}`);

        // 检查问题数据
        let nameUidSame = 0;
        let teamMissing = 0;

        console.log('\n===== 选手数据样例 (每队前3个) =====');
        const shownTeams = new Set();
        players.forEach(p => {
            const nameMatch = p.player_name === p.uid ? '⚠️相同' : '✓不同';
            if (p.player_name === p.uid) nameUidSame++;
            if (!p.team_name) teamMissing++;

            if (!shownTeams.has(p.team_id) && shownTeams.size < 5) {
                console.log(`\n战队: ${p.team_name} (ID=${p.team_id})`);
                console.log(`  选手名="${p.player_name}" UID="${p.uid}" [${nameMatch}]`);
                shownTeams.add(p.team_id);
            }
        });

        // 统计
        console.log('\n===== 数据统计 =====');
        console.log(`姓名与UID相同: ${nameUidSame} 条`);
        console.log(`战队关联缺失: ${teamMissing} 条`);

        // 按战队统计
        const [teamStats] = await connection.execute(`
            SELECT t.team_name, COUNT(p.id) as count
            FROM exa_teams t
            LEFT JOIN exa_players p ON t.id = p.team_id
            GROUP BY t.id, t.team_name
            ORDER BY t.id
        `);

        console.log('\n===== 各战队选手数量 =====');
        teamStats.forEach(t => {
            console.log(`  ${t.team_name}: ${t.count}人`);
        });

    } catch (err) {
        console.error('错误:', err.message);
    } finally {
        if (connection) {
            await connection.end();
        }
    }
}

verifyData();
