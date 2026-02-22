package main

import (
	"fmt"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func staticHandler(w http.ResponseWriter, r *http.Request, conf *Config, logger log.Logger) {
	if conf.Debug {
		_ = level.Debug(logger).Log("task", "Generating static HTML page")
	}

	sessions, err := getAllOpenVPNSessions(conf.ServerName, conf, logger)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OpenVPN Local Clients</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            padding: 40px;
            max-width: 1400px;
            margin: 0 auto;
        }
        h1 {
            color: #333;
            font-size: 2.5em;
            margin-bottom: 30px;
            text-align: center;
        }
        .back-link {
            display: inline-block;
            margin-bottom: 20px;
            color: #667eea;
            text-decoration: none;
            font-weight: 500;
        }
        .back-link:hover {
            text-decoration: underline;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        th {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 15px;
            text-align: left;
            font-weight: 600;
        }
        td {
            padding: 12px 15px;
            border-bottom: 1px solid #eee;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 0.85em;
            font-weight: 500;
        }
        .badge-tcp {
            background-color: #e3f2fd;
            color: #1976d2;
        }
        .badge-udp {
            background-color: #f3e5f5;
            color: #7b1fa2;
        }
        .badge-established {
            background-color: #e8f5e9;
            color: #388e3c;
        }
        .no-data {
            text-align: center;
            padding: 40px;
            color: #999;
            font-size: 1.1em;
        }
        @media (max-width: 768px) {
            .container {
                padding: 20px;
            }
            h1 {
                font-size: 1.8em;
            }
            table {
                font-size: 0.9em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← Back to Dashboard</a>
        <h1>🔒 Active OpenVPN Clients</h1>`

	if len(sessions) == 0 {
		html += `
        <div class="no-data">No active connections</div>`
	} else {
		html += `
        <table>
            <thead>
                <tr>
                    <th>Client ID</th>
                    <th>Remote Host</th>
                    <th>Virtual IP</th>
                    <th>Protocol</th>
                    <th>State</th>
                    <th>Bytes In</th>
                    <th>Bytes Out</th>
                    <th>Connected Since</th>
                </tr>
            </thead>
            <tbody>`

		for _, sess := range sessions {
			protocolClass := "badge-tcp"
			if sess.Protocol == "udp" {
				protocolClass = "badge-udp"
			}

			html += fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%s:%s</td>
                    <td>%s</td>
                    <td><span class="badge %s">%s</span></td>
                    <td><span class="badge badge-established">%s</span></td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`,
				sess.RemoteID,
				sess.RemoteHost,
				sess.RemotePort,
				sess.RemoteTs,
				protocolClass,
				sess.Protocol,
				sess.State,
				sess.BytesIn,
				sess.BytesOut,
				sess.Established,
			)
		}

		html += `
            </tbody>
        </table>`
	}

	html += `
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}