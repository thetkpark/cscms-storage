exports.formatFileSize = (size) => {
    const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    let i = 0;
    while (size >= 1024) {
        size /= 1024;
        ++i;
    }
    return `${Math.round(size * 100) / 100} ${units[i]}`;
}