<script setup lang="ts">
import { Github, Heart, MessageCircle } from 'lucide-vue-next'
</script>

<template>
  <footer class="footer">
    <div class="container footer__container">
      <div class="footer__content">
        <div class="footer__brand">
          <h3>YSM 管理器</h3>
          <p>你的 Minecraft Yes Steve Model 模型库</p>
        </div>
        
        <div class="footer__links">
          <div class="footer__column">
            <h4>链接</h4>
            <a href="/" class="footer__link">首页</a>
            <a href="/upload" class="footer__link">上传模型</a>
          </div>
          
          <div class="footer__column">
            <h4>关于</h4>
            <a href="https://github.com/leemwood/ysmmc" target="_blank" class="footer__link">
              <Github :size="16" /> GitHub
            </a>
          </div>

          <div class="footer__column">
            <h4>社区</h4>
            <a href="https://qm.qq.com/q/fLG2Sc3yU2" target="_blank" class="footer__link" title="YSM 免费模型站用户群">
              <img src="/Tencent_QQ-Logo.wine.svg" alt="QQ" class="icon-img" /> YSM 用户群
            </a>
            <a href="https://qm.qq.com/q/Gjs81js8YC" target="_blank" class="footer__link" title="Fold Craft Launcher玩家交流群①">
              <img src="/Tencent_QQ-Logo.wine.svg" alt="QQ" class="icon-img" /> FCL 玩家群
            </a>
          </div>
        </div>
      </div>
      
      <div class="footer__bottom">
        <p>&copy; {{ new Date().getFullYear() }} YSM 管理器. All rights reserved.</p>
        <p class="footer__credit">
          Made with <Heart :size="14" class="heart-icon" /> by Community
        </p>
      </div>
    </div>
  </footer>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.footer {
  background-color: white;
  border-top: 1px solid var(--color-border);
  padding: $spacing-2xl 0 $spacing-xl;
  margin-top: auto;

  &__container {
    display: flex;
    flex-direction: column;
    gap: $spacing-xl;
  }

  &__content {
    display: flex;
    flex-direction: column;
    gap: $spacing-xl;
    
    @media (min-width: 768px) {
      flex-direction: row;
      justify-content: space-between;
    }
  }

  &__brand {
    h3 {
      font-size: 1.25rem;
      font-weight: 700;
      color: var(--color-text-main);
      margin-bottom: $spacing-xs;
    }
    
    p {
      color: var(--color-text-muted);
      font-size: 0.875rem;
    }
  }

  &__links {
    display: flex;
    gap: $spacing-2xl;
  }

  &__column {
    display: flex;
    flex-direction: column;
    gap: $spacing-sm;
    
    h4 {
      font-size: 0.875rem;
      font-weight: 600;
      color: var(--color-text-main);
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
  }

  &__link {
    color: var(--color-text-muted);
    font-size: 0.875rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    transition: color 0.2s;
    
    &:hover {
      color: var(--color-primary);
    }
  }

  &__bottom {
    border-top: 1px solid var(--color-border);
    padding-top: $spacing-lg;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: $spacing-sm;
    color: var(--color-text-muted);
    font-size: 0.875rem;
    
    @media (min-width: 768px) {
      flex-direction: row;
      justify-content: space-between;
    }
  }

  &__credit {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}

.heart-icon {
  color: var(--color-danger);
  fill: var(--color-danger);
}

.icon-img {
  width: 16px;
  height: 16px;
  object-fit: contain;
}
</style>
