package server

import (
	"auth-server/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

func (s *Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	decoder := json.NewDecoder(r.Body)
	var tokenRequest models.TokenRequest
	err := decoder.Decode(&tokenRequest)
	if err != nil {
		s.HandleError(w, http.StatusBadRequest, TOKEN_ROUTE, err)
	}
	if tokenRequest.GrantType == "client_credentials" {
		ctx := context.Background()
		clientData, appData, err := s.FetchClientAndApplication(ctx, tokenRequest.ClientId, tokenRequest.Aud)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, TOKEN_ROUTE, err)
			return
		}
		s.logger.WithField("clientData", clientData)
		s.logger.WithField("appData", appData)

		if !clientData.HasAllowedScopes(tokenRequest.Scope, appData.AppName) {
			s.HandleError(w, http.StatusUnauthorized, TOKEN_ROUTE, errors.New("client does not have the requested scopes"))
			return
		}

		payload := models.NewPayload(clientData.ID.String(), appData.ID.String(), 1, tokenRequest.Scope)
		jwt, err := models.NewJwt(payload, "JWT")
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, TOKEN_ROUTE, err)
			return
		}
		var tokenResponse models.TokenResponse
		jwtToken, err := jwt.Token()
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, TOKEN_ROUTE, err)
			return
		}
		tokenResponse.AccessToken = jwtToken
		tokenResponse.TokenType = "Bearer"
		response, err := json.Marshal(tokenResponse)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, TOKEN_ROUTE, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		s.logger.Info(http.StatusOK, TOKEN_ROUTE, start)
	}
}

func (s *Server) HandleIntrospection(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) HandleTokenInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) FetchClientAndApplication(ctx context.Context, clientId string, applicationId string) (*models.Client, *models.Application, error) {
	client, err := s.clientRepository.FindById(ctx, clientId)
	if err != nil {
		return nil, nil, err
	}

	application, err := s.applicationRepository.FindById(ctx, applicationId)
	if err != nil {
		return nil, nil, err
	}

	return client, application, nil
}

func (s *Server) ValidateToken(token string) (*models.Payload, error) {
	parts := models.SplitToken(token)
	header, err := models.ParseHeader(parts[0])
	if err != nil {
		return nil, err
	}
	s.logger.WithField("header", header)
	if header.Alg != os.Getenv("AUTH_SERVER_JWT_ALG") {
		return nil, errors.New("invalid algorithm")
	}
	payload, err := models.ParsePayload(parts[1])
	if err != nil {
		return nil, err
	}
	s.logger.WithField("payload", payload)
	if payload.Exp < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	jwt, err := models.NewJwt(payload, "JWT")
	if err != nil {
		return nil, err
	}
	valid, err := jwt.CheckSignature(parts[2])
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid signature")
	}
	return payload, nil
}
