export function parseLogs(text) {
    const rawLogs = text.split('\n').reduce((a, line) => {
        if (a.length === 0) {
            a.push([line]);
        } else {
            const lastLog = a[a.length - 1];
            const lastLogLine = lastLog[lastLog.length - 1];
            if (isEndLogLine(lastLogLine)) {
                a.push([line]);
            } else {
                lastLog.push(line);
            }
        }
        return a;
    }, []);

    const logs = rawLogs.map(lines => {
        if (isEndLogLine(lines[lines.length - 1])) {
            const info = lines.slice(0, lines.length - 2).join('\n');
            const httpObj = parseLine(lines[lines.length - 1]);
            httpObj.info = [info, httpObj.info].filter(i => i).join('\n');
            return httpObj;
        } else {
            return { info: lines.join('\n') };
        }
    });

    console.log(logs);
    return logs;
}
    
export function isEndLogLine (line) {
    return line.indexOf('time=') !== -1 ||
        line.indexOf('[GIN-debug]') !== -1;
}

export function parseLine(line) {
    return (line.indexOf('time=') !== -1) 
        ? parseHTTP(line) 
        : { info: line};
}

export function parseHTTP(line) {
    const idx = line.indexOf('=');
    if (idx === -1) {
        return {};
    } else {
        const name = line.substring(0, idx);
        let value = '';
        let startOfValue = idx + 1;
        let endOfValue = idx + 1;
        if (startOfValue < line.length - 1) {
            if (line[startOfValue] === '"') {
                startOfValue++;
                endOfValue = line.indexOf('"', idx + 2);
                value = line.substring(startOfValue, endOfValue);
                endOfValue++;
            } else {
                endOfValue = line.indexOf(' ', idx + 1);
                value = line.substring(startOfValue, endOfValue);
            }
        }
        const ret = parseHTTP(line.substring(endOfValue + 1));
        ret[name] = value;
        return ret;
    }
}
