// const API_URL = "https://undeliberatingly-decemviral-petronila.ngrok-free.dev/api/v1/user/device-info";
const API_URL = "http://172.16.17.123:8000/api/v1/user/device-info";
const ctx = document.getElementById("chart").getContext("2d");

// Konfigurasi Animasi Smooth
const UPDATE_INTERVAL = 3000; 
const ANIMATION_DURATION = 1100; 

const chart = new Chart(ctx, {
    type: "line",
    data: {
        labels: Array(20).fill(""), 
        datasets: [
            {
                label: "CPU %",
                data: Array(20).fill(0),
                borderColor: "#22c55e",
                backgroundColor: "rgba(34, 197, 94, 0.1)",
                tension: 0.4,
                fill: true,
                pointRadius: 0
            },
            {
                label: "RAM %",
                data: Array(20).fill(0),
                borderColor: "#facc15",
                backgroundColor: "rgba(250, 204, 21, 0.1)",
                tension: 0.4,
                fill: true,
                pointRadius: 0
            },
            {
                label: "Disk %",
                data: Array(20).fill(0),
                borderColor: "#a855f7",
                backgroundColor: "rgba(168, 85, 247, 0.1)",
                tension: 0.4,
                fill: true,
                pointRadius: 0
            }
        ]
    },
    options: {
        responsive: true,
        maintainAspectRatio: false,
        animation: { duration: ANIMATION_DURATION, easing: 'linear' },
        layout: { padding: { right: 10 } },
        scales: {
            y: { min: 0, max: 100, grid: { color: "rgba(255, 255, 255, 0.05)" }, ticks: { color: "#9ca3af" } },
            x: { display: false }
        },
        plugins: { legend: { labels: { color: "#9ca3af" } } }
    }
});

// Fungsi untuk ambil angka di depan % (termasuk yang di dalam kurung)
function parsePercent(str) {
    if (!str) return 0;
    const match = str.match(/(\d+(\.\d+)?)%/);
    return match ? parseFloat(match[1]) : 0;
}

async function loadData() {
    try {
        const res = await fetch(API_URL);
        const json = await res.json();
        const { server, client } = json.data;

        // --- 1. UPDATE CARD UTAMA ---
        document.getElementById("status").textContent = "Live Monitoring";
        document.getElementById("cpuUsage").textContent = server.server_cpu_usage;
        document.getElementById("ramUsage").textContent = server.server_ram_usage;
        document.getElementById("diskUsage").textContent = server.server_disk_usage;
        document.getElementById("healthScore").textContent = server.health_score + "/100";

        // --- 2. UPDATE SERVER INFO (Mapping JSON ke ID) ---
        const serverMap = {
            "serverName": server.server_host_name,
            "serverOS": server.server_os,
            "cpuModel": server.server_cpu_model,
            "uptime": server.system_uptime,
            "appUptime": server.app_uptime,
            "publicIP": server.server_public_ip,
            "internetStatus": server.internet_status,
            "environment": server.environment,
            "suggestion": server.suggestion
        };

        // --- 3. UPDATE CLIENT INFO (Mapping JSON ke ID) ---
        const clientMap = {
            "clientIP": client.user_ip,
            "country": client.country,
            "city": client.city,
            "isp": client.isp,
            "device": `${client.device_brand} ${client.device_model}`.trim(),
            "browser": `${client.browser} ${client.browser_version}`,
            "os": client.os,
            "bot": client.is_bot ? "Yes" : "No",
            "fingerprint": client.fingerprint
        };

        // Fungsi Helper untuk mengisi text jika ID ada di HTML
        const fillText = (map) => {
            for (const [id, val] of Object.entries(map)) {
                const el = document.getElementById(id);
                if (el) el.textContent = val || "-";
            }
        };

        fillText(serverMap);
        fillText(clientMap);

        // --- 4. UPDATE GRAPH ---
        const cpuV = parsePercent(server.server_cpu_usage);
        const ramV = parsePercent(server.server_ram_usage);
        const diskV = parsePercent(server.server_disk_usage);

        chart.data.datasets[0].data.shift();
        chart.data.datasets[0].data.push(cpuV);
        chart.data.datasets[1].data.shift();
        chart.data.datasets[1].data.push(ramV);
        chart.data.datasets[2].data.shift();
        chart.data.datasets[2].data.push(diskV);

        chart.update();

    } catch (e) {
        document.getElementById("status").textContent = "Disconnected";
        console.error("Error fetching data:", e);
    }
}

loadData();
setInterval(loadData, UPDATE_INTERVAL);