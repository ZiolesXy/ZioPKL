// scriptStatic.js
document.addEventListener("DOMContentLoaded", () => {
  const API_URL = "https://undeliberatingly-decemviral-petronila.ngrok-free.dev/api/v1/user/device-info";
  // const API_URL = "http://172.16.17.123:8000/api/v1/user/device-info";
  // const API_URL = "http://localhost:8000/api/v1/user/device-info";
  const UPDATE_INTERVAL = 3000; // 3 detik

  const ctx = document.getElementById("chart").getContext("2d");

  // --- Cache terakhir jika fetch gagal ---
  let lastServer = {};
  let lastClient = {};

  // --- Setup Chart ---
  const chart = new Chart(ctx, {
    type: "line",
    data: {
      labels: Array(10).fill(""),
      datasets: [
        {
          label: "CPU %",
          data: Array(10).fill(0),
          borderColor: "#22c55e",
          backgroundColor: "rgba(34,197,94,0.1)",
          tension: 0.3,
          fill: true,
          pointRadius: 2
        },
        {
          label: "RAM %",
          data: Array(10).fill(0),
          borderColor: "#facc15",
          backgroundColor: "rgba(250,204,21,0.1)",
          tension: 0.3,
          fill: true,
          pointRadius: 2
        },
        {
          label: "Disk %",
          data: Array(10).fill(0),
          borderColor: "#a855f7",
          backgroundColor: "rgba(168,85,247,0.1)",
          tension: 0.3,
          fill: true,
          pointRadius: 2
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      animation: false,
      transitions: { active: { animation: { duration: 1000 } } },
      plugins: { legend: { labels: { color: "#9ca3af" } } },
      scales: {
        y: { min: 0, max: 100, grid: { color: "rgba(255,255,255,0.05)" }, ticks: { color: "#9ca3af" } },
        x: { grid: { display: false }, ticks: { display: false } }
      }
    }
  });

  // --- Helper DOM ---
  function setText(id, value) {
    const el = document.getElementById(id);
    if (el) el.textContent = value ?? "-";
  }

  // --- Load data from API ---
  async function loadData() {
    try {
      const res = await fetch(API_URL, {
        headers: {
          "Accept": "application/json",
          "ngrok-skip-browser-warning": "true"
        }
      });
      const json = await res.json();
      const server = json.data.server;
      const client = json.data.client;

      lastServer = server;
      lastClient = client;

      // --- Update cards ---
      setText("status", "Live Monitoring");
      setText("cpuUsage", server.server_cpu_usage);
      setText("ramUsage", server.server_ram_usage);
      setText("diskUsage", server.server_disk_usage);
      setText("healthScore", `${server.health_score}/100`);

      // --- Server info ---
      setText("serverName", server.server_host_name);
      setText("serverOS", server.server_os);
      setText("cpuModel", server.server_cpu_model);
      setText("uptime", server.system_uptime);
      setText("appUptime", server.app_uptime);
      setText("publicIP", server.server_public_ip);
      setText("internetStatus", server.internet_status);
      setText("environment", server.environment);
      setText("suggestion", server.suggestion);

      // --- Client info ---
      setText("clientIP", client.user_ip);
      setText("device", `${client.device_brand || ""} ${client.device_model || ""}`);
      setText("browser", client.browser);
      setText("os", client.os);
      setText("fingerprint", client.fingerprint);
      setText("country", client.country);
      setText("city", client.city);
      setText("isp", client.isp);
      setText("bot", client.is_bot ? "Yes" : "No");

      // --- Chart update ---
      // Gunakan angka langsung jika backend mengirim server_cpu_percent dsb
      const cpu = parseFloat(server.server_cpu_percent ?? server.server_cpu_usage.replace('%','')) || 0;
      const ramMatch = server.server_ram_percent ?? (server.server_ram_usage.match(/(\d+(\.\d+)?)%/)?.[1]);
      const ram = parseFloat(ramMatch) || 0;
      const diskMatch = server.server_disk_percent ?? (server.server_disk_usage.match(/(\d+(\.\d+)?)%/)?.[1]);
      const disk = parseFloat(diskMatch) || 0;

      chart.data.datasets[0].data.shift(); chart.data.datasets[0].data.push(cpu);
      chart.data.datasets[1].data.shift(); chart.data.datasets[1].data.push(ram);
      chart.data.datasets[2].data.shift(); chart.data.datasets[2].data.push(disk);
      chart.update("active");

    } catch (e) {
      console.error("Fetch Error:", e);
      setText("status", "Disconnected");

      // fallback
      if (lastServer && lastClient) {
        setText("cpuUsage", lastServer.server_cpu_usage);
        setText("ramUsage", lastServer.server_ram_usage);
        setText("diskUsage", lastServer.server_disk_usage);
        setText("healthScore", `${lastServer.health_score}/100`);
        setText("internetStatus", lastServer.internet_status);
      }
    }
  }

  loadData();
  setInterval(loadData, UPDATE_INTERVAL);
});