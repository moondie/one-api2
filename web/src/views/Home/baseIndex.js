import { Box, Typography, Button, Container, Stack } from '@mui/material';
import Grid from '@mui/material/Unstable_Grid2';
import { useSelector } from 'react-redux';
import { useEffect, useState } from 'react';
import { API } from '../../utils/api';
import { showError } from '../../utils/common';

const BaseIndex = () => {
  const account = useSelector((state) => state.account);
  const [token, setToken] = useState('');
  const fetchData = async () => {
    if (!localStorage.getItem('first_apikey')) {
      try {
        const res = await API.get(`/api/token/`, {
          params: {
            page: 1,
            size: 10,
            keyword: '',
            order: '-id'
          }
        });
        const { success, message, data } = res.data;
        if (success) {
          localStorage.setItem('first_apikey', data.data[0].key);
          setToken(data.data[0].key);
        } else {
          showError(message);
        }
      } catch (error) {
        console.error(error);
      }
    } else {
      setToken(localStorage.getItem('first_apikey'));
    }
  };

  useEffect(() => {
    fetchData();
  }, []);
  const server_host = window.location.host;
  return (
    <Box
      sx={{
        minHeight: 'calc(100vh - 136px)',
        backgroundImage: 'linear-gradient(to right, #ff9966, #ff5e62)',
        color: 'white',
        p: 4
      }}
    >
      <Container maxWidth="lg">
        <Grid container columns={12} alignItems="center" sx={{ minHeight: 'calc(100vh - 230px)' }}>
          <Grid md={6} lg={5}>
            <Stack spacing={3}>
              <Typography variant="h1" sx={{ fontSize: '4rem', color: '#fff', lineHeight: 1.5 }}>
                Hust One API
              </Typography>
              <Typography variant="h4" sx={{ fontSize: '1.5rem', color: '#fff', lineHeight: 1.5 }}>
                All in one 的 OpenAI 接口 <br />
                整合各种 API 访问方式 <br />
              </Typography>
              {!account.user ? (
                <Typography>
                  <Button
                    variant="contained"
                    href="/login"
                    sx={{
                      backgroundColor: '#24292e',
                      color: '#fff',
                      width: '122px',
                      height: '50px',
                      boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .3)'
                    }}
                  >
                    登录
                  </Button>
                  <Button
                    variant="contained"
                    href="/register"
                    sx={{
                      backgroundColor: '#fff',
                      color: '#24292e',
                      width: '60px',
                      height: '50px',
                      boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .3)',
                      marginLeft: '20px'
                    }}
                  >
                    注册
                  </Button>
                </Typography>
              ) : (
                <Typography>
                  <Button
                    variant="contained"
                    href={`https://www.hustgpt.com/#/?settings={"key":"sk-${token}","url":"${server_host}"}`}
                    sx={{
                      backgroundColor: 'rgba(134,55,0,0.76)',
                      color: '#fff',
                      width: '180px',
                      height: '50px',
                      boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .3)'
                    }}
                  >
                    去聊天
                  </Button>
                </Typography>
              )}
            </Stack>
          </Grid>

          <Grid md={8} lg={7}>
            <Stack spacing={3}>
              <Typography variant="h1" sx={{ fontSize: '4rem', color: '#985c68', lineHeight: 1.5 }}>
                强力AI助手
              </Typography>
              <Typography variant="h4" sx={{ fontSize: '1.5rem', color: '#3f3e3e', lineHeight: 1.5 }}>
                <div style={{ padding: '5px', marginTop: '8px', backgroundColor: 'rgba(239,224,64,0.58)', borderRadius: '10px' }}>
                  ① 电商: 营销文案、演讲文稿、将自己从繁杂的文字中解放! <br />
                </div>
                <div style={{ padding: '5px', marginTop: '8px', backgroundColor: 'rgba(41,234,245,0.6)', borderRadius: '10px' }}>
                  ② 编程: 帮写代码、修改代码、寻找bug、针对代码特定段落提问，再也不用一次打开十几个*SDN或者博*园窗口啦! <br />
                </div>
                <div style={{ padding: '5px', marginTop: '8px', backgroundColor: 'rgba(239,224,64,0.58)', borderRadius: '10px' }}>
                  ③ 科研: 帮写markdown、论文润色、知识问答，提高科研效率 <br />
                </div>
                <div style={{ padding: '5px', marginTop: '8px', backgroundColor: 'rgba(41,234,245,0.6)', borderRadius: '10px' }}>
                  ④ 生活: 家庭医生、情感导师、健身教练、营养师，成为健康生活达人! <br />
                </div>
              </Typography>
            </Stack>
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

export default BaseIndex;
