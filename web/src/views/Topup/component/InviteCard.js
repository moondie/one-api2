import { Stack, Typography, Container, Box, OutlinedInput, InputAdornment, Button } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import SubCard from 'ui-component/cards/SubCard';
import inviteImage from 'assets/images/invite/cwok_casual_19.webp';
import { useState } from 'react';
import { API } from 'utils/api';
import { showError, copy } from 'utils/common';

const InviteCard = () => {
  const theme = useTheme();
  const [inviteUl, setInviteUrl] = useState('');

  const handleInviteUrl = async () => {
    if (inviteUl) {
      copy(inviteUl, '邀请链接');
      return;
    }

    try {
      const res = await API.get('/api/user/aff');
      const { success, message, data } = res.data;
      if (success) {
        let link = `${window.location.origin}/register?aff=${data}`;
        setInviteUrl(link);
        copy(link, '邀请链接');
      } else {
        showError(message);
      }
    } catch (error) {
      return;
    }
  };

  return (
    <Box component="div">
      <SubCard
        sx={{
          background: theme.palette.primary.dark
        }}
      >
        <Stack justifyContent="center" alignItems={'flex-start'} padding={'40px 24px 0px'} spacing={3}>
          <Container sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <img src={inviteImage} alt="invite" width={'250px'} />
          </Container>
        </Stack>
      </SubCard>
      <SubCard
        sx={{
          marginTop: '-20px'
        }}
      >
        <Stack justifyContent="center" alignItems={'center'} spacing={3}>
          <Typography variant="h3" sx={{ color: theme.palette.primary.dark }}>
            邀请奖励
          </Typography>
          <Typography variant="body" sx={{ color: theme.palette.primary.dark }}>
            新注册用户有1000tokens的试用额度，使用邀请码注册，邀请者和被邀请者额外获得获得20000tokens的试用额度!刷账号将永久封禁邮箱和IP!
          </Typography>

          <OutlinedInput
            id="invite-url"
            label="邀请链接"
            type="text"
            value={inviteUl}
            name="invite-url"
            placeholder="点击生成邀请链接"
            endAdornment={
              <InputAdornment position="end">
                <Button variant="contained" onClick={handleInviteUrl}>
                  {inviteUl ? '复制' : '生成'}
                </Button>
              </InputAdornment>
            }
            aria-describedby="helper-text-channel-quota-label"
            disabled={true}
          />
        </Stack>
      </SubCard>
    </Box>
  );
};

export default InviteCard;
