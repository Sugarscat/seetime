import navbar from './nav/navbar'
import sidebar from './nav/sidebar'

export default {
    base: '/',
    lang: 'zh-CN',
    title: '',
    description: '',
    themeConfig: {
        nav: navbar,
        logo: '/favicon.ico',
        sidebar: sidebar,
        editLink: {
            pattern: 'https://github.com/Sugarscat/hi-automation/edit/master/docs/:path',
            text: '在 GitHub 上编辑此页'
        },
        socialLinks: [
            {icon: 'github', link: 'https://github.com/Sugarscat/hi-automation'}
        ],
        // footer: {
        //     message: 'GNU GENERAL PUBLIC LICENSE V3 Licensed',
        //     copyright: '<a href="https://beian.miit.gov.cn/" target="_blank"></a>'
        // }
    },
    lastUpdated: true,
}
