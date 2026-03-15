const API_URL = "http://172.16.17.123:8000/api/v1/user/device-info";
const ctx = document.getElementById("chart").getContext("2d");

const UPDATE_INTERVAL = 3000; 
const ANIMATION_DURATION = 1000; 
const DATA_COUNT = 20; 

const chart = new Chart(ctx, {
    type: "line",
    data: {
        // KUNCI 1: Tambahkan 1 label ekstra di akhir untuk memberi ruang napas
        labels: Array(DATA_COUNT + 1).fill(""), 
        datasets: [
            {
                label: "CPU %",
                data: Array(DATA_COUNT).fill(null),
                borderColor: "#22c55e",
                backgroundColor: "rgba(34, 197, 94, 0.1)",
                tension: 0.4,
                fill: true,
                pointRadius: 0,
                borderWidth: 2
            },
            {
                label: "RAM %",
                data: Array(DATA_COUNT).fill(null),
                borderColor: "#facc15",
                backgroundColor: "rgba(250, 204, 21, 0.1)",
                tension: 0.4,
                fill: true,
                pointRadius: 0,
                borderWidth: 2
            },
            {
                label: "Disk %",
                data: Array(DATA_COUNT).fill(null),
                borderColor: "#a855f7",
                backgroundColor: "rgba(168, 85, 247, 0.1)",
                tension: 0.4,
                fill: true,
                pointRadius: 0,
                borderWidth: 2
            }
        ]
    },
    options: {
        responsive: true,
        maintainAspectRatio: false,
        // Matikan clipping agar garis digambar utuh sampai ujung
        clip: false, 
        animation: {
            duration: ANIMATION_DURATION,
            easing: 'linear'
        },
        layout: {
            padding: {
                left: -1, // Menutup celah kiri
                right: 0, // Biarkan nol karena kita pakai trik label ekstra
                top: 10
            }
        },
        scales: {
            y: { 
                min: 0, 
                max: 100,
                grid: { color: "rgba(255, 255, 255, 0.05)" },
                ticks: { color: "#9ca3af", padding: 8 }
            },
            x: {
                display: false,
                // KUNCI 2: Paksa sumbu X berhenti di label terakhir yang kosong
                bounds: 'ticks' 
            }
        },
        plugins: {
            legend: { labels: { color: "#9ca3af" } }
        }
    }
});

function parsePercent(str) {
    if (!str) return 0;
    const match = str.match(/(\d+(\.\d+)?)%/);
    return match ? parseFloat(match[1]) : 0;
}

async function loadData() {
    try {
        const res = await fetch(API_URL, { headers: { "ngrok-skip-browser-warning": "true" } });
        const json = await res.json();
        const { server, client } = json.data;

        // Update UI Text
        document.getElementById("status").textContent = "Live Monitoring";
        document.getElementById("cpuUsage").textContent = server.server_cpu_usage;
        document.getElementById("ramUsage").textContent = server.server_ram_usage;
        document.getElementById("diskUsage").textContent = server.server_disk_usage;
        document.getElementById("healthScore").textContent = server.health_score + "/100";

        // Update Detail Info
        const infoMap = {
            "serverName": server.server_host_name,
            "serverOS": server.server_os,
            "cpuModel": server.server_cpu_model,
            "uptime": server.system_uptime,
            "appUptime": server.app_uptime,
            "publicIP": server.server_public_ip,
            "internetStatus": server.internet_status,
            "environment": server.environment,
            "suggestion": server.suggestion,
            "clientIP": client.user_ip,
            "country": client.country,
            "city": client.city,
            "isp": client.isp,
            "device": `${client.device_brand || ''} ${client.device_model || ''}`.trim(),
            "browser": `${client.browser} ${client.browser_version}`,
            "os": client.os,
            "bot": client.is_bot ? "Yes" : "No",
            "fingerprint": client.fingerprint
        };

        for (const [id, val] of Object.entries(infoMap)) {
            const el = document.getElementById(id);
            if (el) el.textContent = val || "-";
        }

        // --- UPDATE DATA GRAPH ---
        const cpuV = parsePercent(server.server_cpu_usage);
        const ramV = parsePercent(server.server_ram_usage);
        const diskV = parsePercent(server.server_disk_usage);

        [cpuV, ramV, diskV].forEach((val, i) => {
            chart.data.datasets[i].data.shift();
            chart.data.datasets[i].data.push(val);
        });

        chart.update();

    } catch (e) {
        document.getElementById("status").textContent = "Disconnected";
        console.error(e);
    }
}

loadData();
setInterval(loadData, UPDATE_INTERVAL);