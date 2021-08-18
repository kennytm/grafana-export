import { minify } from 'terser';
import { createHash } from 'crypto';
import fs from 'fs';
import escapeHTML from 'escape-html';

const terserConfig = {
    toplevel: true,
    ecma: 2015,
    format: {
        max_line_len: 120,
        semicolons: false,
        preamble: '/* Copyright 2020 PingCAP, Inc. Licensed under MIT. */',
    },
    mangle: {
        properties: {
            keep_quoted: true,
            reserved: [
                'annotations',
                'dashboard',
                'expandRows',
                'forEachPanel',
                'getSaveModelClone',
                'hide',
                'injector',
                'on',
                'one',
                'prop',
                'removeClass',
                'removePanel',
                'setAutoRefresh',
                'show',
                'snapshot',
                'snapshotData',
                'startRefresh',
                'timeRange',
            ],
        },
    },
};

async function main() {
    // compute the hash value.
    const content = fs.readFileSync('dist/index.js', {encoding: 'utf-8'});
    const hasher = createHash('sha256');
    hasher.update(content);
    const hash = hasher.digest('hex').substr(0, 12);

    // perform optimization.
    let minified = {};
    for (let language of ['en-US', 'zh-CN', 'zh-TW']) {
        const result = await minify(content, {
            compress: {
                unsafe: true,
                unsafe_arrows: true,
                unsafe_comps: true,
                unsafe_math: true,
                unsafe_undefined: true,
                global_defs: {
                    LANGUAGE: language,
                },
            },
            ...terserConfig
        });
        minified[language] = escapeHTML(result.code);
    }

    // rewrite the website index.html
    let template = fs.readFileSync('../index.html', {encoding: 'utf-8'});
    template = template.replace(
        /<textarea lang="([^"]+)"([^>]*)>.*?<\/textarea>/gs,
        (_, language, attrs) => `<textarea lang="${language}"${attrs}>${minified[language]}</textarea>`
    );
    fs.writeFileSync('../index.html', template);
}

main()
