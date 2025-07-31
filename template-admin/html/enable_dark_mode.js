// 永続ダークモード強制適用スクリプト
console.log('永続ダークモードを強制適用中...');

// ブラウザのlocalStorageでdark-modeを常にtrueに固定
if (typeof localStorage !== 'undefined') {
    localStorage.setItem('dark-mode', 'true');
    // ライトモードへの切り替えを防止
    const originalSetItem = localStorage.setItem;
    localStorage.setItem = function(key, value) {
        if (key === 'dark-mode' && value !== 'true') {
            console.log('ライトモードへの切り替えをブロックしました');
            return originalSetItem.call(this, key, 'true');
        }
        return originalSetItem.call(this, key, value);
    };
}

// HTMLのクラスを強制的にdarkに固定
if (typeof document !== 'undefined') {
    document.documentElement.classList.add('dark');
    document.querySelector('html').style.colorScheme = 'dark';
    
    // DOMの変更を監視してダークモードを維持
    const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
            if (mutation.type === 'attributes' && mutation.attributeName === 'class') {
                const target = mutation.target;
                if (target === document.documentElement && !target.classList.contains('dark')) {
                    target.classList.add('dark');
                    console.log('ダークモードクラスを復元しました');
                }
            }
        });
    });
    
    observer.observe(document.documentElement, {
        attributes: true,
        attributeFilter: ['class']
    });
}

console.log('永続ダークモードが有効になりました - ライトモードは無効化されています');