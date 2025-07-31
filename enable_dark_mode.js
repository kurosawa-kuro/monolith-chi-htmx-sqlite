// ダークモードを有効にするスクリプト
console.log('ダークモードを有効にします...');

// ブラウザのlocalStorageでdark-modeをtrueに設定
if (typeof localStorage !== 'undefined') {
    localStorage.setItem('dark-mode', 'true');
    console.log('localStorage に dark-mode = true を設定しました');
}

// HTMLのクラスを追加してダークモードを適用
if (typeof document !== 'undefined') {
    document.documentElement.classList.add('dark');
    document.querySelector('html').style.colorScheme = 'dark';
    console.log('HTMLにdarkクラスを追加し、カラースキームをdarkに設定しました');
}

console.log('ダークモードが有効になりました！');